// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jomei/notionapi"
)

// Ensure YoloProvider satisfies various provider interfaces.
var _ provider.Provider = &YoloProvider{}

// YoloProvider defines the provider implementation.
type YoloProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// YoloProviderModel describes the provider data model.
type YoloProviderModel struct {
	NotionSecret types.String `tfsdk:"notion_secret"`
}

func (p *YoloProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "yoloexp"
	resp.Version = p.version
}

func (p *YoloProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"notion_secret": schema.StringAttribute{
				MarkdownDescription: "Notion integration secret",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *YoloProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data YoloProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.NotionSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("notion_secret"),
			"Unknown Notion secret.",
			"The provider cannot connect to Notion without a secret.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	notionSecret := os.Getenv("NOTION_SECRET")
	if !data.NotionSecret.IsNull() {
		notionSecret = data.NotionSecret.ValueString()
	}

	if notionSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("notion_secret"),
			"Notion secret is missing.",
			"The provider cannot connect to Notion without a secret.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a notion client.
	client := notionapi.NewClient(
		notionapi.Token(notionSecret),
		notionapi.WithRetry(3),
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *YoloProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *YoloProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &YoloProvider{
			version: version,
		}
	}
}

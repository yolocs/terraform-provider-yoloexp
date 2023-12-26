package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jomei/notionapi"
)

var (
	_ datasource.DataSource              = &notionPageDataSource{}
	_ datasource.DataSourceWithConfigure = &notionPageDataSource{}
)

type notionPageModel struct {
	ID          types.String `tfsdk:"id"`
	URL         types.String `tfsdk:"url"`
	ParentID    types.String `tfsdk:"parent_id"`
	CreatedTime types.String `tfsdk:"created_time"`
}

func NewNotionPageDataSource() datasource.DataSource {
	return &notionPageDataSource{}
}

type notionPageDataSource struct {
	client *notionapi.Client
}

func (d *notionPageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notion_page"
}

func (d *notionPageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Notion page data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Notion page id",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Notion page url",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The page's parent id.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				MarkdownDescription: "The timestamp when this page was created.",
				Computed:            true,
			},
		},
	}
}

func (d *notionPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config notionPageModel
	if err := req.Config.Get(ctx, &config); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	page, err := d.client.Page.Get(ctx, notionapi.PageID(config.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get page",
			fmt.Sprintf("Failed to get page: %s", err),
		)
		return
	}

	state := &notionPageModel{
		ID:          types.StringValue(page.ID.String()),
		URL:         types.StringValue(page.URL),
		ParentID:    types.StringValue(page.Parent.PageID.String()),
		CreatedTime: types.StringValue(page.CreatedTime.String()),
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *notionPageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*notionapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *notionapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

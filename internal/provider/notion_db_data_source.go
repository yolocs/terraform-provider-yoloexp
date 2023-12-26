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
	_ datasource.DataSource              = &notionDatabaseDataSource{}
	_ datasource.DataSourceWithConfigure = &notionDatabaseDataSource{}
)

type notionDatabaseModel struct {
	ID          types.String                  `tfsdk:"id"`
	URL         types.String                  `tfsdk:"url"`
	ParentID    types.String                  `tfsdk:"parent_id"`
	CreatedTime types.String                  `tfsdk:"created_time"`
	Properties  []notionDatabasePropertyModel `tfsdk:"properties"`
}

type notionDatabasePropertyModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func NewNotionDatabaseDataSource() datasource.DataSource {
	return &notionDatabaseDataSource{}
}

type notionDatabaseDataSource struct {
	client *notionapi.Client
}

func (d *notionDatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notion_database"
}

func (d *notionDatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Notion database data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Notion database id",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Notion database url",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The database's parent id.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				MarkdownDescription: "The timestamp when this database was created.",
				Computed:            true,
			},
			"properties": schema.ListNestedAttribute{
				MarkdownDescription: "The properties of the database.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The property's name.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The property's type.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *notionDatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config notionDatabaseModel
	if err := req.Config.Get(ctx, &config); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	db, err := d.client.Database.Get(ctx, notionapi.DatabaseID(config.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get database",
			fmt.Sprintf("Failed to get database: %s", err),
		)
		return
	}

	state := &notionDatabaseModel{
		ID:          types.StringValue(db.ID.String()),
		URL:         types.StringValue(db.URL),
		ParentID:    types.StringValue(db.Parent.PageID.String()),
		CreatedTime: types.StringValue(db.CreatedTime.String()),
	}
	for k, v := range db.Properties {
		state.Properties = append(state.Properties, notionDatabasePropertyModel{
			Name: types.StringValue(k),
			Type: types.StringValue(string(v.GetType())),
		})
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (d *notionDatabaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

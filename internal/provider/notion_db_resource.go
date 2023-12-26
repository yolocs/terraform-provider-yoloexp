package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jomei/notionapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &notionDatabaseResource{}
)

type notionDatabaseResourceModel struct {
	ID          types.String `tfsdk:"id"`
	URL         types.String `tfsdk:"url"`
	ParentID    types.String `tfsdk:"parent_id"`
	CreatedTime types.String `tfsdk:"created_time"`
}

// NewNotionDatabaseResource is a helper function to simplify the provider implementation.
func NewNotionDatabaseResource() resource.Resource {
	return &notionDatabaseResource{}
}

// notionDatabaseResource is the resource implementation.
type notionDatabaseResource struct {
	client *notionapi.Client
}

// Metadata returns the resource type name.
func (r *notionDatabaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notion_database"
}

// Schema defines the schema for the resource.
func (r *notionDatabaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Notion database data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Notion database id",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Notion database url",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The database's parent id. Only supports page_id.",
				Required:            true,
			},
			"created_time": schema.StringAttribute{
				MarkdownDescription: "The timestamp when this database was created.",
				Computed:            true,
			},
			// "properties": schema.ListNestedAttribute{
			// 	MarkdownDescription: "The properties of the database.",
			// 	Computed:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"name": schema.StringAttribute{
			// 				MarkdownDescription: "The property's name.",
			// 				Computed:            true,
			// 			},
			// 			"type": schema.StringAttribute{
			// 				MarkdownDescription: "The property's type.",
			// 				Computed:            true,
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *notionDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notionDatabaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbReq := &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(plan.ParentID.ValueString()),
		},
		Title: []notionapi.RichText{
			{
				Type: notionapi.ObjectTypeText,
				Text: &notionapi.Text{Content: "Title"},
			},
		},
		Properties: notionapi.PropertyConfigs{
			"create": notionapi.TitlePropertyConfig{
				Type: notionapi.PropertyConfigTypeTitle,
			},
		},
		IsInline: false,
	}

	db, err := r.client.Database.Create(ctx, dbReq)
	if err != nil {
		tflog.Debug(ctx, "Failed to create database")
		resp.Diagnostics.AddError(
			"Failed to create database",
			fmt.Sprintf("Failed to create database: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(db.ID.String())
	plan.CreatedTime = types.StringValue(db.CreatedTime.String())
	plan.URL = types.StringValue(db.URL)
	// for k, v := range db.Properties {
	// 	plan.Properties = append(plan.Properties, notionDatabasePropertyModel{
	// 		Name: types.StringValue(k),
	// 		Type: types.StringValue(string(v.GetType())),
	// 	})
	// }

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *notionDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notionDatabaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	db, err := r.client.Database.Get(ctx, notionapi.DatabaseID(state.ID.ValueString()))
	if err != nil {
		tflog.Debug(ctx, "Failed to read database")
		resp.Diagnostics.AddError(
			"Failed to read database",
			fmt.Sprintf("Failed to get database: %s", err),
		)
		return
	}

	state.ID = types.StringValue(db.ID.String())
	state.URL = types.StringValue(db.URL)
	state.ParentID = types.StringValue(db.Parent.PageID.String())
	state.CreatedTime = types.StringValue(db.CreatedTime.String())
	// for k, v := range db.Properties {
	// 	state.Properties = append(state.Properties, notionDatabasePropertyModel{
	// 		Name: types.StringValue(k),
	// 		Type: types.StringValue(string(v.GetType())),
	// 	})
	// }

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notionDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notionDatabaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	db, err := r.client.Database.Get(ctx, notionapi.DatabaseID(plan.ID.ValueString()))
	if err != nil {
		// Assume it's not found. BAD!
		// Only make sure the database exists.
		dbReq := &notionapi.DatabaseCreateRequest{
			Parent: notionapi.Parent{
				Type:   notionapi.ParentTypePageID,
				PageID: notionapi.PageID(plan.ParentID.ValueString()),
			},
			Title: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{Content: "Title"},
				},
			},
			Properties: notionapi.PropertyConfigs{
				"create": notionapi.TitlePropertyConfig{
					Type: notionapi.PropertyConfigTypeTitle,
				},
			},
			IsInline: false,
		}

		db, err = r.client.Database.Create(ctx, dbReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create database",
				fmt.Sprintf("Failed to create database: %s", err),
			)
			return
		}
	}

	plan.ID = types.StringValue(db.ID.String())
	plan.URL = types.StringValue(db.URL)
	plan.ParentID = types.StringValue(db.Parent.PageID.String())
	plan.CreatedTime = types.StringValue(db.CreatedTime.String())
	// for k, v := range db.Properties {
	// 	plan.Properties = append(plan.Properties, notionDatabasePropertyModel{
	// 		Name: types.StringValue(k),
	// 		Type: types.StringValue(string(v.GetType())),
	// 	})
	// }

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *notionDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Not supported.
}

func (d *notionDatabaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Duo Group resource.",

		CreateContext: ResourceGroupCreate,
		ReadContext:   ResourceGroupRead,
		UpdateContext: ResourceGroupUpdate,
		DeleteContext: ResourceGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"desc": {
				Description: "The description of the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"status": {
				Description: "The authentication status of the group. Must be one of: `active` `bypass` `disabled`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Active",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					if val != "Active" && val != "Bypass" && val != "Disabled" {
						errs = append(errs, fmt.Errorf("%q must be one of: 'active', 'bypass', 'disabled', got: %s", key, val))
					}
					return
				},
			},
		},
	}
}

func ResourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	values := url.Values{}
	values.Set("name", d.Get("name").(string))

	if v, ok := d.GetOk("desc"); ok {
		values.Set("desc", v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		values.Set("status", v.(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/groups", values, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	result := &admin.GetGroupResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to create group: %s, error: %s", result.Stat, *result.Message)
	}

	group := result.Response
	d.SetId(group.GroupID)
	tflog.Trace(ctx, "Successfully created group")

	return ResourceGroupRead(ctx, d, meta)
}

func ResourceGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	group_id := d.Id()

	result, err := duoAdminClient.GetGroup(group_id)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Unable to read group: %s, error: %s", result.Stat, *result.Message)
	}

	group := result.Response
	d.Set("name", group.Name)
	d.Set("desc", group.Desc)
	d.Set("status", group.Status)

	return nil
}

func ResourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	group_id := d.Id()
	values := url.Values{}

	d.Partial(true)

	if d.HasChange("name") {
		values.Set("name", d.Get("name").(string))
	}

	if d.HasChange("desc") {
		values.Set("desc", d.Get("desc").(string))
	}

	if d.HasChange("status") {
		values.Set("status", d.Get("status").(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/groups/%s", group_id), values, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	result := &admin.GetGroupResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to update group: %s, error: %s", group_id, *result.Message)
	}
	d.Partial(false)

	return ResourceGroupRead(ctx, d, meta)
}

func ResourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	group_id := d.Id()
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/groups/%s", group_id), nil, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	result := &admin.StringResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to delete group: %s, error: %s", group_id, *result.Message)
	}
	return nil
}

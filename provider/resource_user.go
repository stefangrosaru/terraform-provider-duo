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

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Duo User resource.",

		CreateContext: ResourceUserCreate,
		ReadContext:   ResourceUserRead,
		UpdateContext: ResourceUserUpdate,
		DeleteContext: ResourceUserDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"username": {
				Description: "The name of the user to create.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"realname": {
				Description: "The real name (or full name) of this user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"email": {
				Description: "The email address of this user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"status": {
				Description: "The user's status. Must be one of: `active` `bypass` `disabled`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "active",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					if val != "active" && val != "bypass" && val != "disabled" {
						errs = append(errs, fmt.Errorf("%q must be one of: 'active', 'bypass', 'disabled', got: %s", key, val))
					}
					return
				},
			},
			"notes": {
				Description: "An optional description or notes field. Can be viewed in the Duo Admin Panel.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"firstname": {
				Description: "The user's given name.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"lastname": {
				Description: "The user's surname.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func ResourceUserCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	values := url.Values{}
	values.Set("username", d.Get("username").(string))

	if v, ok := d.GetOk("realname"); ok {
		values.Set("realname", v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		values.Set("email", v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		values.Set("status", v.(string))
	}

	if v, ok := d.GetOk("notes"); ok {
		values.Set("notes", v.(string))
	}

	if v, ok := d.GetOk("firstname"); ok {
		values.Set("firstname", v.(string))
	}

	if v, ok := d.GetOk("lastname"); ok {
		values.Set("lastname", v.(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/users", values, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	result := &admin.GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to create user: %s, error: %s", result.Stat, *result.Message)
	}

	user := result.Response
	d.SetId(user.UserID)
	tflog.Trace(ctx, "Successfully created user")

	return ResourceUserRead(ctx, d, meta)
}

func ResourceUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	user_id := d.Id()

	result, err := duoAdminClient.GetUser(user_id)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Unable to read user: %s, error: %s", result.Stat, *result.Message)
	}

	user := result.Response
	d.Set("username", user.Username)
	d.Set("realname", user.RealName)
	d.Set("email", user.Email)
	d.Set("status", user.Status)
	d.Set("notes", user.Notes)
	d.Set("firstname", user.FirstName)
	d.Set("lastname", user.LastName)

	return nil
}

func ResourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	user_id := d.Id()
	values := url.Values{}

	d.Partial(true)

	if d.HasChange("username") {
		values.Set("username", d.Get("username").(string))
	}

	if d.HasChange("realname") {
		values.Set("realname", d.Get("realname").(string))
	}

	if d.HasChange("email") {
		values.Set("email", d.Get("email").(string))
	}

	if d.HasChange("status") {
		values.Set("status", d.Get("status").(string))
	}

	if d.HasChange("notes") {
		values.Set("notes", d.Get("notes").(string))
	}

	if d.HasChange("firstname") {
		values.Set("firstname", d.Get("firstname").(string))
	}

	if d.HasChange("lastname") {
		values.Set("lastname", d.Get("lastname").(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/users/%s", user_id), values, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	result := &admin.GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to update user: %s, error: %s", user_id, *result.Message)
	}
	d.Partial(false)

	return ResourceUserRead(ctx, d, meta)
}

func ResourceUserDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	user_id := d.Id()
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/users/%s", user_id), nil, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	result := &admin.StringResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to delete user: %s, error: %s", user_id, *result.Message)
	}
	return nil
}

package internal

import (
	"context"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Provides details about a specific Duo User.",

		ReadContext: DataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "The ID of the user to retrieve.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": {
				Description: "The name of the user to retrieve.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"realname": {
				Description: "The real name (or full name) of this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The email address of this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The user's status. Must be one of: `active` `bypass` `disabled`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"notes": {
				Description: "An optional description or notes field. Can be viewed in the Duo Admin Panel.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"firstname": {
				Description: "The user's given name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"lastname": {
				Description: "The user's surname.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func DataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	user_id := d.Get("user_id").(string)

	result, err := duoAdminClient.GetUser(user_id)

	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			return diag.Errorf("Unable to read user: %s, error: %s", result.Stat, *result.Message)
		}
		return diag.Errorf("Unable to read user: %s, error: %s", result.Stat, *result.Message)
	}

	user := result.Response

	d.SetId(user.UserID)
	d.Set("username", user.Username)
	d.Set("realname", user.RealName)
	d.Set("email", user.Email)
	d.Set("status", user.Status)
	d.Set("notes", user.Notes)
	d.Set("firstname", user.FirstName)
	d.Set("lastname", user.LastName)

	return nil
}

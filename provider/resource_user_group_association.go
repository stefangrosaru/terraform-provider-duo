package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUserGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Duo Group resource.",

		CreateContext: ResourceUserGroupAssociationCreate,
		ReadContext:   ResourceUserGroupAssociationRead,
		DeleteContext: ResourceUserGroupAssociationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the group to associate with the user.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"user_id": {
				Description: "The ID of the user to associate with the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func ResourceUserGroupAssociationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	group_id := d.Get("group_id").(string)
	user_id := d.Get("user_id").(string)

	values := url.Values{}
	values.Set("group_id", group_id)

	_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/users/%s/groups", user_id), values, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	result := &admin.StringResult{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to add user to group: %s, error: %s", result.Stat, *result.Message)
	}

	id := strings.Join([]string{group_id, user_id}, "-")
	d.SetId(id)
	tflog.Trace(ctx, "Successfully added user to group")

	return ResourceUserGroupAssociationRead(ctx, d, meta)
}

func ResourceUserGroupAssociationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := d.Id()

	s := strings.Split(id, "-")
	group_id, user_id := s[0], s[1]

	d.Set("group_id", group_id)
	d.Set("user_id", user_id)

	return nil
}

func ResourceUserGroupAssociationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	id := d.Id()

	s := strings.Split(id, "-")
	group_id, user_id := s[0], s[1]

	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/users/%s/groups/%s", user_id, group_id), nil, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}

	result := &admin.StringResult{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return diag.Errorf("An error has occurred: %s", err)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to remove user from group: %s, error: %s", group_id, *result.Message)
	}

	return nil
}

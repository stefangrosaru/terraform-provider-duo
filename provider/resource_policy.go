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

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Duo Group resource.",

		CreateContext: ResourcePolicyCreate,
		ReadContext:   ResourcePolicyRead,
		UpdateContext: ResourcePolicyUpdate,
		DeleteContext: ResourcePolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"new_user_policy": {
				Description: "Enable policy for new users.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enroll_policy": {
				Description: "enroll_policy",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"auth_status_activated": {
				Description: "auth_status_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"auth_status": {
				Description: "auth_status",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"user_locations_activated": {
				Description: "user_locations_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"user_locations_default_action": {
				Description: "user_locations_default_action",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint_health_activated": {
				Description: "endpoint_health_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint_health_policy_macos": {
				Description: "endpoint_health_policy_macos",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint_health_enroll_policy_macos": {
				Description: "endpoint_health_enroll_policy_macos",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint_health_policy_windows": {
				Description: "endpoint_health_policy_windows",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint_health_enroll_policy_windows": {
				Description: "endpoint_health_enroll_policy_windows",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"trusted_sessions_activated": {
				Description: "trusted_sessions_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"trusted_devices_activated": {
				Description: "trusted_devices_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"platforms_activated": {
				Description: "platforms_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"android_allowed": {
				Description: "android_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"android_warn_policy_version": {
				Description: "android_warn_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"android_block_policy_version": {
				Description: "android_block_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"blackberry_allowed": {
				Description: "blackberry_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"chrome_os_allowed": {
				Description: "chrome_os_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ios_allowed": {
				Description: "ios_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ios_warn_policy_version": {
				Description: "ios_warn_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ios_block_policy_version": {
				Description: "ios_block_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"linux_allowed": {
				Description: "linux_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"macos_allowed": {
				Description: "macos_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"macos_warn_policy_version": {
				Description: "macos_warn_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"macos_block_policy_version": {
				Description: "macos_block_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"windows_allowed": {
				Description: "windows_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"windows_warn_policy_version": {
				Description: "windows_warn_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"windows_block_policy_version": {
				Description: "windows_block_policy_version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"windows_phone_allowed": {
				Description: "windows_phone_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"other_os_allowed": {
				Description: "other_os_allowed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"browsers_activated": {
				Description: "browsers_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"plugins_activated": {
				Description: "plugins_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"flash_remediation": {
				Description: "flash_remediation",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"java_remediation": {
				Description: "java_remediation",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"networks_activated": {
				Description: "networks_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"networks_allow": {
				Description: "networks_allow",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"networks_2fa": {
				Description: "networks_2fa",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"anonymous_ip_policy": {
				Description: "anonymous_ip_policy",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"anonymous_ip_policy_activated": {
				Description: "anonymous_ip_policy_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"risk_based_factor_selection_activated": {
				Description: "risk_based_factor_selection_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"auth_methods_activated": {
				Description: "auth_methods_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_factor_push": {
				Description: "allow_factor_push",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_factor_mobile_otp": {
				Description: "allow_factor_mobile_otp",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_factor_sms": {
				Description: "allow_factor_sms",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_factor_web_auth": {
				Description: "allow_factor_web_auth",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"web_auth_policies": {
				Description: "web_auth_policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_factor_hard_token": {
				Description: "allow_factor_hard_token",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"duo_mobile_app_activated": {
				Description: "duo_mobile_app_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"require_updated_duo_mobile": {
				Description: "require_updated_duo_mobile",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mobile_rooted_devices_activated": {
				Description: "mobile_rooted_devices_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_rooted_devices": {
				Description: "allow_rooted_devices",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mobile_lock_activated": {
				Description: "mobile_lock_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"require_lock": {
				Description: "require_lock",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mobile_encryption_activated": {
				Description: "mobile_encryption_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"require_encryption": {
				Description: "require_encryption",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mobile_touch_id_activated": {
				Description: "mobile_touch_id_activated",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"require_touch_id": {
				Description: "require_touch_id",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func ResourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	duoClient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoClient)

	values := url.Values{}
	values.Set("name", d.Get("name").(string))

	if v, ok := d.GetOk("new_user_policy"); ok {
		values.Set("new-user-policy-activated", v.(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", "/policies", values, duoapi.UseTimeout)
	if err != nil {
		return diag.Errorf("1An error has occurred: %s", err)
	}

	result := &admin.StringResult{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return diag.Errorf("2An error has occurred: %s", result.Response)
	}
	if result.Stat != "OK" {
		return diag.Errorf("Unable to add user to group: %s, error: %s", result.Stat, *result.Message)
	}

	tflog.Trace(ctx, "Successfully added user to group")

	return ResourcePolicyRead(ctx, d, meta)
}

func ResourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// id := d.Id()

	// s := strings.Split(id, "-")
	// group_id, user_id := s[0], s[1]

	// d.Set("group_id", group_id)
	// d.Set("user_id", user_id)

	return nil
}

func ResourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// duoClient := meta.(*duoapi.DuoApi)
	// duoAdminClient := admin.New(*duoClient)

	// user_id := d.Id()
	// values := url.Values{}

	// d.Partial(true)

	// if d.HasChange("username") {
	// 	values.Set("username", d.Get("username").(string))
	// }

	// if d.HasChange("realname") {
	// 	values.Set("realname", d.Get("realname").(string))
	// }

	// if d.HasChange("email") {
	// 	values.Set("email", d.Get("email").(string))
	// }

	// if d.HasChange("status") {
	// 	values.Set("status", d.Get("status").(string))
	// }

	// if d.HasChange("notes") {
	// 	values.Set("notes", d.Get("notes").(string))
	// }

	// if d.HasChange("firstname") {
	// 	values.Set("firstname", d.Get("firstname").(string))
	// }

	// if d.HasChange("lastname") {
	// 	values.Set("lastname", d.Get("lastname").(string))
	// }

	// _, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/users/%s", user_id), values, duoapi.UseTimeout)
	// if err != nil {
	// 	return diag.Errorf("An error has occurred: %s", err)
	// }

	// result := &admin.GetUserResult{}
	// err = json.Unmarshal(body, result)
	// if err != nil {
	// 	return diag.Errorf("An error has occurred: %s", err)
	// }
	// if result.Stat != "OK" {
	// 	return diag.Errorf("Unable to update user: %s, error: %s", user_id, *result.Message)
	// }
	// d.Partial(false)

	return ResourceUserRead(ctx, d, meta)
}

func ResourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

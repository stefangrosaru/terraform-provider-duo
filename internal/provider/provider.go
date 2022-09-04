package provider

import (
	"context"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stefangrosaru/terraform-provider-duo/internal/user"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"integration_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("DUO_INTEGRATION_KEY", nil),
					Description: "Duo Admin API Integration key",
				},
				"secret_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("DUO_SECRET_KEY", nil),
					Description: "Duo Admin API Secret skey",
				},
				"api_hostname": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("DUO_API_HOSTNAME", nil),
					Description: "Duo Admin API Server hostname",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"duo_user": user.DataSourceUser(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"duo_user": user.ResourceUser(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// Setup a User-Agent for the API client
		user_agent := "terraform-provider-duo/" + version

		// Get API client credential
		integration_key := d.Get("integration_key").(string)
		secret_key := d.Get("secret_key").(string)
		api_hostname := d.Get("api_hostname").(string)

		client := duoapi.NewDuoApi(integration_key, secret_key, api_hostname, user_agent)

		return client, nil
	}
}

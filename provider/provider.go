package provider

import (
	"os"

	"github.com/comwork/comwork-provider/api/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("PROVIDER_REGION"); v != "" {
					  return v, nil
					}
			
					return "fr-par-1", nil
				},			
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_TOKEN", ""),
			},
			"ngx_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_NGINX_USERNAME", ""),
			},
			"ngx_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_NGINX_PASSWORD", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"comwork_instance": resourceInstance(),

		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	region := d.Get("region").(string)
	token := d.Get("token").(string)
	ngx_username := d.Get("ngx_username").(string)
	ngx_password := d.Get("ngx_password").(string)

	return client.NewClient(region, token,ngx_username,ngx_password), nil

}

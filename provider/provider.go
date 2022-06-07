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
		},
		ResourcesMap: map[string]*schema.Resource{
			"comwork_instance": resourceInstance(),
			"comwork_project": resourceProject(),

		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	region := d.Get("region").(string)
	token := d.Get("token").(string)

	return client.NewClient(region, token), nil

}

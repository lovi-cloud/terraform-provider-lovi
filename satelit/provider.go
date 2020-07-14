package satelit

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var descriptions map[string]string

// Provider provide schema
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("SATELIT_API_ENDPOINT", ""),
				Description: descriptions["api_endpoint"],
				Required:    true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"satelit_volume":            resourceSatelitVolume(),
			"satelit_volume_attachment": resourceSatelitVolumeAttachment(),
			"satelit_virtual_machine":   resourceSatelitVirtualMachine(),
		},
		ConfigureFunc: configureProvider,
	}
}

func init() {
	descriptions = map[string]string{
		"api_endpoint": "The endpoint of Satelit server",
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIEndpoint: d.Get("api_endpoint").(string),
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

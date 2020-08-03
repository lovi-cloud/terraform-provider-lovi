package lovi

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var descriptions map[string]string

// Provider provide schema
func Provider() *schema.Provider {
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
			"lovi_volume":            resourceLoviVolume(),
			"lovi_volume_attachment": resourceLoviVolumeAttachment(),
			"lovi_virtual_machine":   resourceLoviVirtualMachine(),
			"lovi_subnet":            resourceLoviSubnet(),
			"lovi_bridge":            resourceLoviBridge(),
			"lovi_internal_bridge":   resourceLoviInternalBridge(),
			"lovi_address":           resourceLoviAddress(),
			"lovi_lease":             resourceLoviLease(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func init() {
	descriptions = map[string]string{
		"api_endpoint": "The endpoint of Satelit server",
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := Config{
		APIEndpoint: d.Get("api_endpoint").(string),
	}

	if err := config.LoadAndValidate(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to load and validate config: %v", err),
		})

		return nil, diags
	}

	return &config, diags
}

package lovi

import (
	"context"
	"fmt"

	satelitpb "github.com/whywaita/satelit/api/satelit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceInternalBridge() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInternalBridgeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 10),
			},
		},
	}
}

func dataSourceInternalBridgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	resp, err := client.ListBridge(ctx, &satelitpb.ListBridgeRequest{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call ListBridge (internal): %+v", err),
		})

		return diags
	}

	n := d.Get("name").(string)

	var bridge *satelitpb.Bridge
	for _, b := range resp.Bridges {
		if b.Name == n {
			bridge = b
			break
		}
	}
	if bridge == nil {
		// not found
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("bridge is not found (internal)"),
		})

		return diags
	}

	d.SetId(bridge.Uuid)
	d.Set("name", bridge.Name)

	return diags
}

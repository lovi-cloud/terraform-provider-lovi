package lovi

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	satelitpb "github.com/lovi-cloud/satelit/api/satelit"
)

func dataSourceLoviBridge() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoviBridgeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 10),
				ExactlyOneOf: []string{"name", "vlan_id"},
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"name", "vlan_id"},
			},
		},
	}
}

func dataSourceLoviBridgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	resp, err := client.ListBridge(ctx, &satelitpb.ListBridgeRequest{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call ListBridge: %+v", err),
		})

		return diags
	}

	var bridge *satelitpb.Bridge
	if name, ok := d.GetOk("name"); ok {
		n := name.(string)
		for _, b := range resp.Bridges {
			if b.Name == n {
				bridge = b
				break
			}
		}
	} else if vlanID, ok := d.GetOk("vlan_id"); ok {
		id := vlanID.(int)
		i := uint32(id)

		for _, b := range resp.Bridges {
			if b.VlanId == i {
				bridge = b
				break
			}
		}
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Illegal state: one of name or vlan_id must be set"),
		})

		return diags
	}

	if bridge == nil {
		// not found
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("bridge is not found"),
		})

		return diags
	}

	d.SetId(bridge.Uuid)
	d.Set("name", bridge.Name)
	d.Set("vlan_id", bridge.VlanId)

	return diags
}

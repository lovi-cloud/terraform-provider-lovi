package lovi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	satelitpb "github.com/lovi-cloud/satelit/api/satelit"
)

func resourceLoviBridge() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviBridgeCreate,
		ReadContext:   resourceLoviBridgeRead,
		DeleteContext: resourceLoviBridgeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 10),
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoviBridgeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	vlanIDInt := d.Get("vlan_id").(int)

	req := &satelitpb.CreateBridgeRequest{
		Name:   d.Get("name").(string),
		VlanId: int32(vlanIDInt),
	}

	resp, err := client.CreateBridge(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call CreateBridge: %v", err),
		})

		return diags
	}

	d.SetId(resp.Bridge.Uuid)
	resourceLoviBridgeRead(ctx, d, meta)

	return diags
}

func resourceLoviBridgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.GetBridgeRequest{
		Uuid: d.Id(),
	}
	resp, err := client.GetBridge(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call GetBridge: %v", err),
		})

		return diags
	}

	d.Set("name", resp.Bridge.Name)
	d.Set("vlan_id", resp.Bridge.VlanId)

	return diags
}

func resourceLoviBridgeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteBridgeRequest{
		Uuid: d.Id(),
	}
	_, err := client.DeleteBridge(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteBridge: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

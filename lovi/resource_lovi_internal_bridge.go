package lovi

import (
	"context"
	"fmt"
	"time"

	satelitpb "github.com/whywaita/satelit/api/satelit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLoviInternalBridge() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviInternalBridgeCreate,
		ReadContext:   resourceLoviInternalBridgeRead,
		DeleteContext: resourceLoviInternalBridgeDelete,
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
		},
	}
}

func resourceLoviInternalBridgeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.CreateInternalBridgeRequest{
		Name: d.Get("name").(string),
	}
	resp, err := client.CreateInternalBridge(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call CreateInternalBridge: %v", err),
		})

		return diags
	}

	d.SetId(resp.Bridge.Uuid)
	resourceLoviInternalBridgeRead(ctx, d, meta)

	return diags
}

func resourceLoviInternalBridgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			Summary:  fmt.Sprintf("failed to call GetBridge (internal): %v", err),
		})

		return diags
	}

	d.Set("name", resp.Bridge.Name)

	return diags
}

func resourceLoviInternalBridgeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			Summary:  fmt.Sprintf("failed to call DeleteBridge (internal): %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

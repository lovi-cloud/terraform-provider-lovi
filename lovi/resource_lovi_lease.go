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

func resourceLoviLease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviLeaseCreate,
		ReadContext:   resourceLoviLeaseRead,
		DeleteContext: resourceLoviLeaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"address_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceLoviLeaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.CreateLeaseRequest{
		AddressId: d.Get("address_id").(string),
	}
	resp, err := client.CreateLease(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call CreateLease: %v", err),
		})

		return diags
	}

	d.SetId(resp.Lease.Uuid)
	resourceLoviLeaseRead(ctx, d, meta)

	return diags
}

func resourceLoviLeaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.GetLeaseRequest{
		Uuid: d.Id(),
	}
	resp, err := client.GetLease(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call GetLease: %v", err),
		})

		return diags
	}

	d.Set("mac_address", resp.Lease.MacAddress)

	return diags
}
func resourceLoviLeaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteLeaseRequest{
		Uuid: d.Id(),
	}
	_, err := client.DeleteLease(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteLease: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

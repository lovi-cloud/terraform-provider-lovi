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

func resourceLoviAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviAddressCreate,
		ReadContext:   resourceLoviAddressRead,
		DeleteContext: resourceLoviAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceLoviAddressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.CreateAddressRequest{
		SubnetId: d.Get("subnet_id").(string),
	}
	resp, err := client.CreateAddress(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call CreateAddress: %v", err),
		})

		return diags
	}

	d.SetId(resp.Address.Uuid)
	resourceLoviAddressRead(ctx, d, meta)

	return diags
}

func resourceLoviAddressRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.GetAddressRequest{
		Uuid: d.Id(),
	}
	resp, err := client.GetAddress(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call GetAddress: %v", err),
		})

		return diags
	}

	d.Set("subnet_id", resp.Address.SubnetId)
	d.Set("ip", resp.Address.Ip)

	return diags
}

func resourceLoviAddressDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteAddressRequest{
		Uuid: d.Id(),
	}
	_, err := client.DeleteAddress(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteAddress: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

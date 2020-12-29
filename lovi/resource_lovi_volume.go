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

func resourceLoviVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviVolumeCreate,
		ReadContext:   resourceLoviVolumeRead,
		DeleteContext: resourceLoviVolumeDelete,
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
				ValidateFunc: validation.IsUUID,
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "GB",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"backend_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "europa's backend name",
			},
		},
	}
}

func resourceLoviVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	sizeInt := d.Get("size").(int)

	req := &satelitpb.AddVolumeRequest{
		Name:             d.Get("name").(string),
		CapacityGigabyte: uint32(sizeInt),
		BackendName:      d.Get("backend_name").(string),
	}
	resp, err := client.AddVolume(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call AddVolume: %v", err),
		})
		return diags
	}

	volume := resp.Volume
	d.SetId(volume.Id)
	resourceLoviVolumeRead(ctx, d, meta)

	return diags
}

func resourceLoviVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.ShowVolumeRequest{
		Id: d.Id(),
	}
	resp, err := client.ShowVolume(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call ShowVolume: %v", err),
		})
		return diags
	}

	d.Set("attached", resp.Volume.Attached)
	d.Set("hostname", resp.Volume.Hostname)
	d.Set("capacity_byte", resp.Volume.CapacityGigabyte)
	d.Set("backend_name", resp.Volume.BackendName)

	return diags
}

func resourceLoviVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteVolumeRequest{
		Id: d.Id(),
	}
	_, err := client.DeleteVolume(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteVolume: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

package lovi

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	satelitpb "github.com/whywaita/satelit/api/satelit"
)

func resourceLoviVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviVolumeAttachmentCreate,
		ReadContext:   resourceLoviVolumeAttachmentRead,
		DeleteContext: resourceLoviVolumeAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoviVolumeAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	volumeID := d.Get("volume_id").(string)
	hostname := d.Get("hostname").(string)

	req := &satelitpb.AttachVolumeRequest{
		Id:       volumeID,
		Hostname: hostname,
	}
	_, err := client.AttachVolume(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call AttachVolume: %v", err),
		})

		return diags
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s", volumeID, hostname)))
	resourceLoviVolumeAttachmentRead(ctx, d, meta)

	return diags
}

func resourceLoviVolumeAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	volumeID := d.Get("volume_id").(string)
	hostname := d.Get("hostname").(string)

	req := &satelitpb.ShowVolumeRequest{
		Id: volumeID,
	}
	resp, err := client.ShowVolume(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call ShowVolume: %v", err),
		})

		return diags
		// TODO: check not found (need to delete if not found)
	}
	volume := resp.Volume

	if volume.Attached == false || volume.Hostname == "" || volume.Hostname != hostname {
		d.SetId("")
	}

	return diags
}

func resourceLoviVolumeAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	volumeID := d.Get("volume_id").(string)

	req := &satelitpb.DetachVolumeRequest{
		Id: volumeID,
	}
	_, err := client.DetachVolume(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DetachVolume: %v", err),
		})

		return diags
	}
	d.SetId("")

	return diags
}

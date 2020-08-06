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

func resourceLoviInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviInterfaceAttachmentCreate,
		ReadContext:   resourceLoviInterfaceAttachmentRead,
		DeleteContext: resourceLoviInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"virtual_machine_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"bridge_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"average": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 10),
			},
			"lease_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceLoviInterfaceAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	averageInt := d.Get("average").(int)

	req := &satelitpb.AttachInterfaceRequest{
		VirtualMachineId: d.Get("virtual_machine_id").(string),
		BridgeId:         d.Get("bridge_id").(string),
		Average:          int64(averageInt),
		Name:             d.Get("name").(string),
		LeaseId:          d.Get("lease_id").(string),
	}
	resp, err := client.AttachInterface(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call AttachInterface: %v", err),
		})

		return diags
	}

	d.SetId(resp.InterfaceAttachment.Uuid)
	resourceLoviInterfaceAttachmentRead(ctx, d, meta)

	return diags
}

func resourceLoviInterfaceAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.GetAttachmentRequest{
		AttachmentId: d.Id(),
	}
	resp, err := client.GetAttachment(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call GetAttachment: %v", err),
		})

		return diags
	}

	d.Set("virtual_machine_id", resp.InterfaceAttachment.VirtualMachineId)
	d.Set("bridge_id", resp.InterfaceAttachment.BridgeId)
	d.Set("average", resp.InterfaceAttachment.Average)
	d.Set("name", resp.InterfaceAttachment.Name)
	d.Set("lease_id", resp.InterfaceAttachment.LeaseId)

	return diags
}

func resourceLoviInterfaceAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DetachInterfaceRequest{
		AtttachmentId: d.Id(),
	}
	_, err := client.DetachInterface(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DetachInterface: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

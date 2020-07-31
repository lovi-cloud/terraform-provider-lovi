package lovi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	satelitpb "github.com/whywaita/satelit/api/satelit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLoviVirtualMachine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviVirtualMachineCreate,
		ReadContext:   resourceLoviVirtualMachineRead,
		DeleteContext: resourceLoviVirtualMachineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "virtual machine's display name",
			},
			"vcpus": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"memory_kib": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"root_volume_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"source_image_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "source image uuid",
				ValidateFunc: validation.IsUUID,
			},
			"hypervisor_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoviVirtualMachineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	vcpus := d.Get("vcpus").(int)
	memoryKib := d.Get("memory_kib").(int)
	rootVolumeGB := d.Get("root_volume_gb").(int)
	sourceImageID := d.Get("source_image_id").(string)
	hypervisorName := d.Get("hypervisor_name").(string)

	req := &satelitpb.AddVirtualMachineRequest{
		Name:           name,
		Vcpus:          uint32(vcpus),
		MemoryKib:      uint64(memoryKib),
		RootVolumeGb:   uint32(rootVolumeGB),
		SourceImageId:  sourceImageID,
		HypervisorName: hypervisorName,
	}
	resp, err := client.AddVirtualMachine(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call AddVirtualMachine: %v", err),
		})

		return diags
	}

	vmUUID := resp.Uuid
	d.SetId(vmUUID)
	resourceLoviVirtualMachineRead(ctx, d, meta)

	return diags
}

func resourceLoviVirtualMachineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.ShowVirtualMachineRequest{
		Uuid: d.Id(),
	}
	resp, err := client.ShowVirtualMachine(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call ShowVirtualMachine: %v", err),
		})

		return diags
	}

	d.Set("name", resp.VirtualMachine.Name)
	d.Set("vcpus", resp.VirtualMachine.Vcpus)
	d.Set("memory_kib", resp.VirtualMachine.MemoryKib)
	d.Set("hypervisor_name", resp.VirtualMachine.HypervisorName)

	return diags
}

func resourceLoviVirtualMachineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteVirtualMachineRequest{
		Uuid: d.Id(),
	}
	_, err := client.DeleteVirtualMachine(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteVirtualMachine: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

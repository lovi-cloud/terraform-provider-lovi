package lovi

import (
	"context"
	"time"

	satelitpb "github.com/whywaita/satelit/api/satelit"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceLoviVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoviVirtualMachineCreate,
		Read:   resourceLoviVirtualMachineRead,
		Delete: resourceLoviVirtualMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceLoviVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

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
	resp, err := client.AddVirtualMachine(context.Background(), req)
	if err != nil {
		return err
	}

	vmUUID := resp.Uuid
	d.SetId(vmUUID)

	return resourceLoviVirtualMachineRead(d, meta)
}

func resourceLoviVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	req := &satelitpb.ShowVirtualMachineRequest{
		Uuid: d.Id(),
	}
	resp, err := client.ShowVirtualMachine(context.Background(), req)
	if err != nil {
		return err
	}

	d.Set("name", resp.VirtualMachine.Name)
	d.Set("vcpus", resp.VirtualMachine.Vcpus)
	d.Set("memory_kib", resp.VirtualMachine.MemoryKib)
	d.Set("hypervisor_name", resp.VirtualMachine.HypervisorName)

	return nil
}

func resourceLoviVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	req := &satelitpb.DeleteVirtualMachineRequest{
		Uuid: d.Id(),
	}
	_, err := client.DeleteVirtualMachine(context.Background(), req)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

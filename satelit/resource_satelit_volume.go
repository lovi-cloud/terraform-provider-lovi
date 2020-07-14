package satelit

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	satelitpb "github.com/whywaita/satelit/api/satelit"
)

func resourceSatelitVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceSatelitVolumeCreate,
		Read:   resourceSatelitVolumeRead,
		Delete: resourceSatelitVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}
}

func resourceSatelitVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	sizeInt := d.Get("size").(int)

	req := &satelitpb.AddVolumeRequest{
		Name:             d.Get("name").(string),
		CapacityGigabyte: uint32(sizeInt),
	}
	resp, err := client.AddVolume(context.Background(), req)
	if err != nil {
		return err
	}
	volume := resp.Volume

	d.SetId(volume.Id)

	return resourceSatelitVolumeRead(d, meta)
}

func resourceSatelitVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	req := &satelitpb.ShowVolumeRequest{
		Id: d.Id(),
	}
	resp, err := client.ShowVolume(context.Background(), req)
	if err != nil {
		return err
	}

	d.Set("attached", resp.Volume.Attached)
	d.Set("hostname", resp.Volume.Hostname)
	d.Set("capacity_byte", resp.Volume.CapacityGigabyte)

	return nil
}

func resourceSatelitVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	req := &satelitpb.DeleteVolumeRequest{
		Id: d.Id(),
	}
	_, err := client.DeleteVolume(context.Background(), req)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

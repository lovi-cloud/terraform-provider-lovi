package lovi

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	satelitpb "github.com/whywaita/satelit/api/satelit"
)

func resourceLoviVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoviVolumeAttachmentCreate,
		Read:   resourceLoviVolumeAttachmentRead,
		Delete: resourceLoviVolumeAttachmentDelete,

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

func resourceLoviVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	volumeID := d.Get("volume_id").(string)
	hostname := d.Get("hostname").(string)

	req := &satelitpb.AttachVolumeRequest{
		Id:       volumeID,
		Hostname: hostname,
	}
	_, err := client.AttachVolume(context.Background(), req)
	if err != nil {
		return err
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s", volumeID, hostname)))
	return resourceLoviVolumeAttachmentRead(d, meta)
}

func resourceLoviVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	volumeID := d.Get("volume_id").(string)
	hostname := d.Get("hostname").(string)

	req := &satelitpb.ShowVolumeRequest{
		Id: volumeID,
	}
	resp, err := client.ShowVolume(context.Background(), req)
	if err != nil {
		return err
		// TODO: check not found (need to delete if not found)
	}
	volume := resp.Volume

	if volume.Attached == false || volume.Hostname == "" || volume.Hostname != hostname {
		d.SetId("")
	}

	return nil
}

func resourceLoviVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.SatelitClient

	volumeID := d.Get("volume_id").(string)

	req := &satelitpb.DetachVolumeRequest{
		Id: volumeID,
	}
	_, err := client.DetachVolume(context.Background(), req)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

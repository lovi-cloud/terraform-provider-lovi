package lovi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	satelitpb "github.com/whywaita/satelit/api/satelit"
)

func resourceLoviSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviSubnetCreate,
		ReadContext:   resourceLoviSubnetRead,
		DeleteContext: resourceLoviSubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"start": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"end": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dns_server": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"metadata_server": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourceLoviSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	vlanIDInt := d.Get("vlan_id").(int)

	req := &satelitpb.CreateSubnetRequest{
		Name:           d.Get("name").(string),
		VlanId:         uint32(vlanIDInt),
		Network:        d.Get("network").(string),
		Start:          d.Get("start").(string),
		End:            d.Get("end").(string),
		Gateway:        d.Get("gateway").(string),
		DnsServer:      d.Get("dns_server").(string),
		MetadataServer: d.Get("metadata_server").(string),
	}

	resp, err := client.CreateSubnet(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call CreateSubnet: %v", err),
		})

		return diags
	}

	d.SetId(resp.Subnet.Uuid)
	resourceLoviSubnetRead(ctx, d, meta)

	return diags
}

func resourceLoviSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.GetSubnetRequest{
		Uuid: d.Id(),
	}
	resp, err := client.GetSubnet(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call GetSubnet: %v", err),
		})
	}

	d.Set("name", resp.Subnet.Name)
	d.Set("vlan_id", resp.Subnet.VlanId)
	d.Set("network", resp.Subnet.Network)
	d.Set("start", resp.Subnet.Start)
	d.Set("end", resp.Subnet.End)
	d.Set("gateway", resp.Subnet.Gateway)
	d.Set("dns_server", resp.Subnet.DnsServer)
	d.Set("metadata_server", resp.Subnet.MetadataServer)

	return diags
}

func resourceLoviSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteSubnetRequest{
		Uuid: d.Id(),
	}

	_, err := client.DeleteSubnet(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteSubnet: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

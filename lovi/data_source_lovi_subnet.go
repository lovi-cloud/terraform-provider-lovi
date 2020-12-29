package lovi

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	satelitpb "github.com/lovi-cloud/satelit/api/satelit"
)

func dataSourceLoviSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoviSubnetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLoviSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	var subnet *satelitpb.Subnet

	if id, ok := d.GetOk("id"); ok {
		req := &satelitpb.GetSubnetRequest{
			Uuid: id.(string),
		}
		resp, err := client.GetSubnet(ctx, req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("failed to call GetSubnet: %v", err),
			})

			return diags
		}

		subnet = resp.Subnet
	} else if name, ok := d.GetOk("name"); ok {
		resp, err := client.ListSubnet(ctx, &satelitpb.ListSubnetRequest{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("failed to call ListSubnet: %v", err),
			})

			return diags
		}
		n := name.(string)

		for _, s := range resp.Subnets {
			if s.Name == n {
				subnet = s

				break
			}
		}

		if subnet == nil {
			// not found
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("subnet is not found (name: %v)", n),
			})

			return diags
		}
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Illegal state: one of id or name must be set"),
		})

		return diags
	}

	d.SetId(subnet.Uuid)
	d.Set("name", subnet.Name)
	d.Set("vlan_id", subnet.VlanId)
	d.Set("network", subnet.Network)
	d.Set("start", subnet.Start)
	d.Set("end", subnet.End)
	d.Set("gateway", subnet.Gateway)
	d.Set("dns_server", subnet.DnsServer)
	d.Set("metadata_server", subnet.MetadataServer)

	return diags
}

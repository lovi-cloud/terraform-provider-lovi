package lovi

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	satelitpb "github.com/lovi-cloud/satelit/api/satelit"
)

func resourceLoviCPUPinningGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoviCPUPinningGroupCreate,
		ReadContext:   resourceLoviCPUPinningGroupRead,
		DeleteContext: resourceLoviCPUPinningGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"count_of_core": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"hypervisor_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoviCPUPinningGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	countOFCore := d.Get("count_of_core").(int)
	hypervisorName := d.Get("hypervisor_name").(string)

	req := &satelitpb.AddCPUPinningGroupRequest{
		Name:           name,
		CountOfCore:    uint32(countOFCore),
		HypervisorName: hypervisorName,
	}
	resp, err := client.AddCPUPinningGroup(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call AddCPUPinningGroup: %v", err),
		})

		return diags
	}

	d.SetId(resp.CpuPinningGroup.Uuid)
	resourceLoviCPUPinningGroupRead(ctx, d, meta)

	return diags
}

func resourceLoviCPUPinningGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.ShowCPUPinningGroupRequest{
		Uuid: d.Id(),
	}
	resp, err := client.ShowCPUPinningGroup(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call ShowCPUPinningGroup: %v", err),
		})

		return diags
	}

	d.Set("name", resp.CpuPinningGroup.Name)
	d.Set("count_of_core", resp.CpuPinningGroup.CountOfCore)

	return diags
}

func resourceLoviCPUPinningGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client := config.SatelitClient
	var diags diag.Diagnostics

	req := &satelitpb.DeleteCPUPinningGroupRequest{
		Uuid: d.Id(),
	}
	_, err := client.DeleteCPUPinningGroup(ctx, req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to call DeleteCPUPinningGroup: %v", err),
		})

		return diags
	}

	d.SetId("")
	return diags
}

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/whywaita/terraform-provider-lovi/lovi"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: lovi.Provider})
}

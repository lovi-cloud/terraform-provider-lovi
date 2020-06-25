package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/whywaita/terraform-provider-satelit/satelit"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: satelit.Provider})
}

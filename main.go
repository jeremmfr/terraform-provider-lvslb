package main

import (
	"terraform-provider-lvslb/lvslb"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: lvslb.Provider,
	})
}

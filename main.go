package main

import (
	"github.com/jeremmfr/terraform-provider-lvslb/lvslb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: lvslb.Provider,
	})
}

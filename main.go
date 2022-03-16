package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/amirghedira/terraform-provider/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}

package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: nutanix.Provider,
		})
	}
	KnownPlugins["terraform-provider-nutanix"] = exec
}
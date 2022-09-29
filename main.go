package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/stefangrosaru/terraform-provider-duo/internal"
)

var (
	version string = "dev"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug: debugMode,

		ProviderAddr: "registry.terraform.io/stefangrosaru/duo",

		ProviderFunc: internal.New(version),
	}

	plugin.Serve(opts)
}

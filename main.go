package main

import (
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-routeros/terraform-provider-routeros/routeros"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name routeros

// Set by GoReleaser.
var (
	version = "dev"
	commit  = "none"
)

func main() {
	var debug, showVersion bool

	// https://github.com/hashicorp/terraform-docs-common/blob/main/website/docs/plugin/debugging.mdx
	// https://developer.hashicorp.com/terraform/plugin/debugging
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s (%s)\n", version, commit)
		return
	}

	plugin.Serve(&plugin.ServeOpts{
		ProviderAddr: "registry.terraform.io/shinbunbun/routeros",
		ProviderFunc: routeros.NewProvider,
		Debug:        debug,
	})
}

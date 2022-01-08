package install

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/consul/command/flags"
	"github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/src"
	"github.com/mitchellh/cli"
)

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.init()
	return c
}

type cmd struct {
	UI    cli.Ui
	flags *flag.FlagSet
	http  *flags.HTTPFlags
	help  string

	// flags
	flagSource      bool
	flagDestination bool

	// testStdin is the input for testing.
	testStdin io.Reader
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 2
	}

	args = c.flags.Args()
	if len(args) != 1 {
		c.UI.Error(fmt.Sprintf("Error: command requires exactly one argument: plugin name"))
		return 1
	}

	ctx := context.Background()
	home, err := homedir.Dir()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error: unable to determine home directory: %s", err))
		return 1
	}
	pluginDir := filepath.Join(home, ".consul", "plugins")
	if err := os.MkdirAll(pluginDir, 0700); err != nil {
		c.UI.Error(fmt.Sprintf("Error: unable to create plugin dir: %s", err))
		return 1
	}
	installer := install.NewInstaller()
	execPath, err := installer.Install(ctx, []src.Installable{
		&releases.LatestVersion{
			Product:                  product.Product{
				Name:              fmt.Sprintf("consul-%s", args[0]),
				BinaryName:        func() string {return fmt.Sprintf("consul-%s", args[0])},
				GetVersion:        nil,
				BuildInstructions: nil,
			},
			//Constraints:              nil,
			InstallDir:               pluginDir,
			Timeout:                  0,
			IncludePrereleases:       false,
			SkipChecksumVerification: false,
			ArmoredPublicKey:         "",
		},
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error: unable to install: %s", err))
		return 1
	}
	c.UI.Output(fmt.Sprintf("Installed successfully into %s", execPath))

	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

const (
	synopsis = "Install a plugin."
	help     = `
Usage: consul plugin install NAME [options]

  Install a specific plugin. The plugin will be installed into the
  directory set by the CONSUL_PLUGIN_DIR environment variable that
  defaults to ~/.consul/plugins.

      $ consul plugin install k8s
      $ consul plugin install k8s -version 9.9.9

`
)

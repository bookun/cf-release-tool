package main

import (
	"flag"
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/bookun/cf-release-tool/client"
	"github.com/bookun/cf-release-tool/controller"
	"github.com/bookun/cf-release-tool/manager"
	"github.com/bookun/cf-release-tool/usecase"
)

// Plug has flag information
type Plug struct {
	file   *string
	branch *string
}

// Run is exectuted for the first time
// This Method is implements about Run method in code.cloudfoundry.org/cli/plugin
func (c *Plug) Run(cliConnection plugin.CliConnection, args []string) {

	releaseFlagSet := flag.NewFlagSet("release", flag.ExitOnError)
	manifestFile := releaseFlagSet.String("f", "manifest.yml", "The app will be released based on this manifest file")
	branch := releaseFlagSet.String("b", "master", "An app is released by using this branch")
	if err := releaseFlagSet.Parse(args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//client := client.NewDummyClient(os.Stdout)
	client := client.NewClient(cliConnection)
	manager := manager.NewManager(client)
	inputPort := usecase.NewUsecase(manager)
	ctl := &controller.Controller{
		InputPort:    inputPort,
		InfoGetter:   client,
		ManifestFile: *manifestFile,
		Branch:       *branch,
	}
	if err := ctl.Release(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetMetadata has plugin information
// This Method is implements about GetMetadata method in code.cloudfoundry.org/cli/plugin
func (c *Plug) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "ReleaseTool",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 4,
		},
		Commands: []plugin.Command{
			{
				Name:     "release",
				HelpText: "This plugin executes BlueGreenDeployment for PHP app based on git branch. use --help",
				UsageDetails: plugin.Usage{
					Usage: "release front or tool App\n	cf release [-f] <manifest file>  [-b] <branch>",
					Options: map[string]string{
						"file":   "input manifest file's path",
						"branch": "input git branch name that you will release",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(Plug))
}

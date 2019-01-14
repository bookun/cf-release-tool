package main

import (
	"flag"
	"fmt"
	"github.com/bookun/cf-release-tool/entity"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/bookun/cf-release-tool/client"
	"github.com/bookun/cf-release-tool/controller"
	"github.com/bookun/cf-release-tool/usecase"
)

// Plug has flag information
type Plug struct {
	file   *string
	branch *string
	host   *string
	force  *bool
}

// Run is exectuted for the first time
// This Method is implements about Run method in code.cloudfoundry.org/cli/plugin
func (c *Plug) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	releaseFlagSet := flag.NewFlagSet("release", flag.ExitOnError)
	manifestFile := releaseFlagSet.String("f", "manifest.yml", "The app will be released based on this manifest file")
	branch := releaseFlagSet.String("b", "", "An app is released by using this branch")
	host := releaseFlagSet.String("n", "", "An app is released with hostname")
	force := releaseFlagSet.Bool("y", false, "Answer yes for all question")
	if err := releaseFlagSet.Parse(args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//client := client.NewDummyClient(os.Stdout)
	manifest, err := entity.NewManifest(*manifestFile, *branch, *host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, app := range manifest.Applications {
		client := client.NewClient(cliConnection, app, *force, *manifestFile)
		inputPort := usecase.NewUsecase(app.Name, client)
		ctl := &controller.Controller{
			InputPort:  inputPort,
			InfoGetter: client,
		}
		if err := ctl.Release(); err != nil {
			fmt.Println(err)
			continue
		}
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
					Usage: "release front or tool App\n	cf release [-y] [-f] <manifest file>  [-b] <branch> [-n] <hostname>",
					Options: map[string]string{
						"y":        "answer yes for all questions",
						"file":     "input manifest file's path",
						"branch":   "input git branch name that you will release",
						"hostname": "input hostname",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(Plug))
}

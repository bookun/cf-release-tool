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

type MockPlug struct {
	file   *string
	branch *string
}

// Run はCF plugin では最初に起動されるメソッド
func (c *MockPlug) Run(cliConnection plugin.CliConnection, args []string) {

	releaseFlagSet := flag.NewFlagSet("release", flag.ExitOnError)
	manifestFile := releaseFlagSet.String("f", "manifest.yml", "The app will be released based on this manifest file")
	branch := releaseFlagSet.String("b", "master", "An app is released by using this branch")
	if err := releaseFlagSet.Parse(args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client := client.NewDummyClient(cliConnection)
	manager := manager.NewManager(client)
	inputPort := usecase.NewUsecase(manager)
	ctl := &controller.Controller{
		InputPort:    inputPort,
		ManifestFile: *manifestFile,
		Branch:       *branch,
	}
	if err := ctl.Release(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *MockPlug) GetMetadata() plugin.PluginMetadata {
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
				Alias:    "gootop-11",
				HelpText: "CFアプリのリリースをします。use --help",
				UsageDetails: plugin.Usage{
					Usage: "release front or tool App\n	cf release [-f] <manifest file>  [-b] <branch>",
					Options: map[string]string{
						"file":   "リリース対象のmanifestファイルを選択してください",
						"branch": "リリース対象のブランチを選択してください",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(MockPlug))
}

package client

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
)

type DummyClient struct {
	CC plugin.CliConnection
}

func NewDummyClient(cc plugin.CliConnection) *DummyClient {
	return &DummyClient{
		CC: cc,
	}
}

// Init で bp-configを適切なものに差し替え、
// リリース対象のブランチを最新化。
// 最後にspaceをリリース対象に切り替える
func (c *DummyClient) Init(materialDir, branch, org, space string) error {
	fmt.Printf("cp -r %s ./.bp-config\n", materialDir)
	fmt.Printf("git checkout %s\n", branch)
	fmt.Printf("git pull origin %s\n", branch)
	fmt.Printf("target -o %s -s %s\n", org, space)
	return nil
}

// Push で指定した名前のアプリを cf push
func (c *DummyClient) Push(app, manifestFile string) error {
	//TODO: manifest-front_blue なってまう
	fmt.Printf("push %s -f %s\n", app, manifestFile)
	return nil
}

// Rename で名前の変更を行う
func (c *DummyClient) Rename(oldApp, newApp string) error {
	fmt.Printf("rename %s to %s\n", oldApp, newApp)
	return nil
}

// Delete でAppの削除を行う
func (c *DummyClient) Delete(app string) error {
	fmt.Printf("delete %s\n", app)
	return nil
}

// MapRoute で appにURLをつける
func (c *DummyClient) MapRoute(app, domain, host string) error {
	if host != "" {
		fmt.Printf("map-route %s %s --hostname %s\n", app, domain, host)
	} else {
		fmt.Printf("map-route %s %s\n", app, domain)
	}
	return nil
}

// UnMapRoute で appからURLを取る
func (c *DummyClient) UnMapRoute(app, domain, host string) error {
	if host != "" {
		fmt.Printf("unmap-route %s %s --hostname %s\n", app, domain, host)
	} else {
		fmt.Printf("unmap-route %s %s\n", app, domain)
	}
	return nil
}

func (c *DummyClient) show() {
	fmt.Println("以下の条件でリリースを行います")
}

func (c *DummyClient) confirm(question, ans string) bool {
	fmt.Printf("%s: ", question)
	var userAns string
	fmt.Scan(&userAns)
	if userAns == ans {
		return true
	}
	return false
}

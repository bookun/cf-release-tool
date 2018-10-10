package client

import (
	"fmt"
	"os/exec"

	"github.com/cloudfoundry/cli/plugin"
)

// Client は CFの操作を行うもの
type Client struct {
	cc plugin.CliConnection
}

// NewClient は Clientを初期化
func NewClient(cc plugin.CliConnection) *Client {
	return &Client{
		cc: cc,
	}
}

// Init で bp-configを適切なものに差し替え、
// リリース対象のブランチを最新化。
// 最後にspaceをリリース対象に切り替える
func (c *Client) Init(materialDir, branch, org, space string) error {
	exec.Command("rm", "-rf", "./.bp-config").Run()
	if err := exec.Command("cp", "-rf", materialDir, "./.bp-config").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "checkout", branch).Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "pull", "origin", branch).Run(); err != nil {
		return err
	}
	if _, err := c.cc.CliCommand("target", "-o", org, "-s", space); err != nil {
		return err
	}
	return nil
}

// Push で指定した名前のアプリを cf push
func (c *Client) Push(app, manifestFile string) error {
	if _, err := c.cc.CliCommand("push", app, "-f", manifestFile); err != nil {
		return err
	}
	return nil
}

// Rename で名前の変更を行う
func (c *Client) Rename(oldApp, newApp string) error {
	if _, err := c.cc.CliCommand("rename", oldApp, newApp); err != nil {
		return err
	}
	return nil
}

// Delete でAppの削除を行う
func (c *Client) Delete(app string) error {
	if _, err := c.cc.CliCommand("delete", app); err != nil {
		return err
	}
	return nil
}

// MapRoute で appにURLをつける
func (c *Client) MapRoute(app, domain, host string) error {
	if host != "" {
		if _, err := c.cc.CliCommand("map-route", app, domain, "--hostname", host); err != nil {
			return err
		}
	} else {
		if _, err := c.cc.CliCommand("map-route", app, domain); err != nil {
			return err
		}
	}
	return nil
}

// UnMapRoute で appからURLを取る
func (c *Client) UnMapRoute(app, domain, host string) error {
	if host != "" {
		if _, err := c.cc.CliCommand("unmap-route", app, domain, "--hostname", host); err != nil {
			return err
		}
	} else {
		if _, err := c.cc.CliCommand("unmap-route", app, domain); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) show() {
	fmt.Println("以下の条件でリリースを行います")
}

func (c *Client) confirm(question, ans string) bool {
	fmt.Printf("%s: ", question)
	var userAns string
	fmt.Scan(&userAns)
	if userAns == ans {
		return true
	}
	return false
}

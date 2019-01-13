package client

import (
	"fmt"
	"github.com/bookun/cf-release-tool/entity"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Client operates cloudfoundry API.
type Client struct {
	cc           plugin.CliConnection
	app          entity.App
	force        bool
	manifestFile string
}

// NewClient init Client
func NewClient(cc plugin.CliConnection, app entity.App, force bool) *Client {
	return &Client{
		cc:    cc,
		app:   app,
		force: force,
	}
}

// Init prepare material, git branch, and cf target.
func (c *Client) Init() error {
	for from, to := range c.app.Env.Copy {
		if _, err := os.Stat(to); err == nil {
			if err := exec.Command("rm", "-fr", to).Run(); err != nil {
				err = fmt.Errorf("failed to remove %s before copy", to)
				return err
			}
		}
		if err := exec.Command("cp", "-fr", from, to).Run(); err != nil {
			err = fmt.Errorf("failed to copy from %s to %s", from, to)
			return err
		}
	}
	branch := c.app.Env.Branch
	if branch != "" {
		if err := exec.Command("git", "checkout", branch).Run(); err != nil {
			err = fmt.Errorf("failed to checkout branch")
			return err
		}
		if err := exec.Command("git", "pull", "origin", branch).Run(); err != nil {
			err = fmt.Errorf("failed to pull branch")
			return err
		}
	}
	if _, err := c.cc.CliCommand("target", "-o", c.app.Env.Org, "-s", c.app.Env.Space); err != nil {
		return err
	}
	return nil
}

// Push executes cf push.
func (c *Client) Push(name string) error {
	if _, err := c.cc.CliCommand("push", name, "-f", c.manifestFile); err != nil {
		return err
	}
	return nil
}

// RenameFrom executes cf rename from arg to target app name.
func (c *Client) Rename(from, to string) error {
	if _, err := c.cc.CliCommand("rename", from, to); err != nil {
		return err
	}
	return nil
}

// Stop executes cf stop
func (c *Client) Stop(name string) error {
	if _, err := c.cc.CliCommand("stop", name); err != nil {
		return err
	}
	return nil
}

// Delete executes cf delete
func (c *Client) DeleteApps(appKind string) error {
	apps, err := c.cc.GetApps()
	var appNames []string
	if err != nil {
		return err
	}
	for _, v := range apps {
		if strings.Contains(v.Name, appKind) && v.Name != appKind {
			appNames = append(appNames, v.Name)
		}
	}
	if len(appNames) > 3 {
		sort.Strings(appNames)
		for _, v := range appNames[:len(appNames)-3] {
			if c.force {
				if _, err := c.cc.CliCommand("delete", "-f", v); err != nil {
					return err
				}
			} else {
				if _, err := c.cc.CliCommand("delete", v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// MapRoute executes cf map-route
func (c *Client) MapRoute(name string) error {
	confirmFlag := false
	domain := c.app.Domain
	host := c.app.Host
	if domain == "" {
		err := fmt.Errorf("domain is not be set")
		return err
	}
	if name != c.app.Name {
		host = name
		confirmFlag = true
		if c.app.Env.TestUp != nil {
			domain = c.app.Env.TestUp["domain"]
			host = c.app.Env.TestUp["host"]
		}
	}
	if host == "" {
		if _, err := c.cc.CliCommand("map-route", name, domain); err != nil {
			return err
		}
		return nil
	}
	if _, err := c.cc.CliCommand("map-route", name, domain, "--hostname", host); err != nil {
		return err
	}
	if confirmFlag {
		if err := c.confirm(); err != nil {
			return err
		}
	}
	return nil
}

// UnMapRoute executes cf unmap-route
func (c *Client) UnMapRoute(app string) error {
	appInfo, err := c.cc.GetApp(app)
	if err != nil {
		return err
	}
	for _, route := range appInfo.Routes {
		domain := route.Domain.Name
		host := route.Host
		if domain == "" {
			err := fmt.Errorf("the domain is not be set")
			return err
		}
		if host != "" {
			if _, err := c.cc.CliCommand("unmap-route", app, domain, "--hostname", host); err != nil {
				return err
			}
			if domain != c.app.Domain || host != c.app.Host {
				if _, err := c.cc.CliCommand("delete-route", domain, "--hostname", host); err != nil {
					return err
				}
			}
			return nil
		}
		if _, err := c.cc.CliCommand("unmap-route", app, domain); err != nil {
			return err
		}
		if domain != c.app.Domain || host != c.app.Host {
			if _, err := c.cc.CliCommand("delete-route", domain, "--hostname", host); err != nil {
				return err
			}
		}
	}
	return nil
}

// TestUp execute map-route test host
//func (c *Client) TestUp(app, domain string) (bool, error) {
//	if !c.force {
//		var confirm string
//		tempHost := fmt.Sprintf("test-%s-%s", app, strconv.FormatInt(time.Now().Unix(), 10))
//		domain := c.app.Env.TestUp["domain"]
//		host := c.app.Env.TestUp["host"]
//		if domain == "" || host == "" {
//			//err := fmt.Errorf("domain or host for testup is not be set")
//			return true, nil
//		}
//		if err := c.MapRoute(app, domain, host); err != nil {
//			return false, err
//		}
//		fmt.Printf("Is it displayed properly? [y/n]")
//		if _, err := fmt.Scan(&confirm); err != nil {
//			return false, err
//		}
//		if confirm == "y" {
//			if err := c.UnMapRoute(app); err != nil {
//				return false, err
//			}
//			if err := c.DeleteRoute(domain, tempHost); err != nil {
//				return false, err
//			}
//			return true, nil
//		}
//		if err := c.UnMapRoute(app); err != nil {
//			return false, err
//		}
//		if err := c.Delete(app); err != nil {
//			return false, err
//		}
//		return false, nil
//	}
//	return true, nil
//}

func (c *Client) confirm() error {
	if c.force {
		return nil
	}
	var confirm string
	fmt.Printf("Is it displayed properly? [y/n]")
	if _, err := fmt.Scan(&confirm); err != nil {
		return err
	}
	if confirm != "y" {
		err := fmt.Errorf("deploy is cenceled by you")
		return err
	}
	return nil
}

// CreateBlueName execute naming.
// Blue app is named app name + created time
func (c *Client) CreateBlueName(app string) (string, error) {
	appModel, err := c.cc.GetApp(app)
	if err != nil {
		return "", nil
	}
	timeStr := appModel.PackageUpdatedAt.Format("20060102150405")
	name := fmt.Sprintf("%s-%s", app, timeStr)
	return name, nil
}

// AppExists check if there is a app in your space
func (c *Client) AppExists() error {
	_, err := c.cc.GetApp(c.app.Name)
	if err != nil {
		return err
	}
	return nil
}

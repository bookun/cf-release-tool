package client

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/cli/plugin"
)

// Client operates cloudfoundry API.
type Client struct {
	cc plugin.CliConnection
}

// NewClient init Client
func NewClient(cc plugin.CliConnection) *Client {
	return &Client{
		cc: cc,
	}
}

// Init prepare material, git branch, and cf target.
func (c *Client) Init(envFile, materialDir, branch, org, space string) error {
	if envFile != "" {
		exec.Command("cp", envFile, "./.env").Run()
	}
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

// Push executes cf push.
func (c *Client) Push(app, manifestFile string) error {
	if _, err := c.cc.CliCommand("push", app, "-f", manifestFile); err != nil {
		return err
	}
	return nil
}

// Rename executes cf rename.
func (c *Client) Rename(oldApp, newApp string) error {
	if _, err := c.cc.CliCommand("rename", oldApp, newApp); err != nil {
		return err
	}
	return nil
}

// Delete executes cf delete
func (c *Client) Delete(app string) error {
	apps, err := c.cc.GetApps()
	var appNames []string
	if err != nil {
		return err
	}
	for _, v := range apps {
		if strings.Contains(v.Name, app) && v.Name != app {
			appNames = append(appNames, v.Name)
		}
	}
	if len(appNames) > 3 {
		sort.Strings(appNames)
		for _, v := range appNames[3:] {
			if _, err := c.cc.CliCommand("delete", v); err != nil {
				return err
			}
		}
	}
	return nil
}

// MapRoute executes cf map-route
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

// UnMapRoute executes cf unmap-route
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

// DeleteRoute execute cf delete-route
func (c *Client) DeleteRoute(domain, host string) error {
	if _, err := c.cc.CliCommand("delete-route", "-f", domain, "-n", host); err != nil {
		return err
	}
	return nil
}

// TestUp execute map-route test host
func (c *Client) TestUp(app, domain string) (bool, error) {
	var confirm string
	tempHost := fmt.Sprintf("test-%s-%s", app, strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.MapRoute(app, domain, tempHost); err != nil {
		return false, err
	}
	fmt.Printf("Is it displayed properly? [y/n]")
	fmt.Scan(&confirm)
	if confirm == "y" {
		if err := c.UnMapRoute(app, domain, tempHost); err != nil {
			return false, err
		}
		if err := c.DeleteRoute(domain, tempHost); err != nil {
			return false, err
		}
		return true, nil
	}
	if err := c.UnMapRoute(app, domain, tempHost); err != nil {
		return false, err
	}
	if err := c.Delete(app); err != nil {
		return false, err
	}
	return false, nil
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
func (c *Client) AppExists(app string) error {
	_, err := c.cc.GetApp(app)
	if err != nil {
		return err
	}
	return nil
}

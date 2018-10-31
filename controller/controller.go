package controller

import (
	"errors"
	"io/ioutil"

	"github.com/bookun/cf-release-tool/entity"
	"github.com/bookun/cf-release-tool/usecase"
	yaml "gopkg.in/yaml.v2"
)

// CurrentInfoGetter interface has AppExists method.
// AppExists method is implemented in client package.
type CurrentInfoGetter interface {
	AppExists(app string) error
}

// Controller decides usecase, BlueGreenDeployment or normal Deployment.
type Controller struct {
	InputPort    usecase.InputPort
	InfoGetter   CurrentInfoGetter
	ManifestFile string
	Branch       string
	Host         string
}

// Manifest has information of manifest.yml
type Manifest struct {
	Applications []struct {
		Name       string `yaml:"name"`
		Instance   int    `yaml:"instances"`
		Memory     string `yaml:"momery"`
		Buildpack  string `yaml:"buildpack"`
		NoHostName bool   `yaml:"no-hostname"`
		NoRoute    bool   `yaml:"no-route"`
		Env        struct {
			Org      string `yaml:"ORG"`
			Space    string `yaml:"SPACE"`
			TimeZone string `yaml:"TZ"`
			Lang     string `yaml:"LANG"`
			Host     string `yaml:"HOST"`
			Domain   string `yaml:"DOMAIN"`
			Material string `yaml:"MATERIAL"`
		} `yaml:"env"`
	} `yaml:"applications"`
}

// Release executes deployment.
// If there is application you want to release, this method executes normal deployment.
// Else, it executes BlueGreenDeployment.
func (c *Controller) Release() error {
	m, err := c.getManifest()
	if err != nil {
		return err
	}
	targetApps := m.Applications
	for _, targetApp := range targetApps {
		domain := targetApp.Env.Domain
		host := targetApp.Env.Host
		if len(targetApps) == 1 {
			host, err = c.getHostName()
			if err != nil {
				return err
			}
		}
		entity := entity.Deploy{
			Org:          targetApp.Env.Org,
			Space:        targetApp.Env.Space,
			App:          targetApp.Name,
			ManifestFile: c.ManifestFile,
			MaterialDir:  targetApp.Env.Material,
			Branch:       c.Branch,
		}
		if c.InfoGetter.AppExists(entity.App) != nil {
			if err := c.InputPort.Deployment(entity, domain, host); err != nil {
				return err
			}
		} else {
			if err := c.InputPort.BlueGreenDeployment(entity, domain, host); err != nil {
				return err
			}
		}
	}
	return nil
}

// getManifest parse manifest file.
func (c *Controller) getManifest() (Manifest, error) {
	data, err := ioutil.ReadFile(c.ManifestFile)
	if err != nil {
		return Manifest{}, err
	}
	m := Manifest{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func (c *Controller) getHostName() (string, error) {
	if c.Host != "" {
		return c.Host, nil
	}
	err := errors.New("host name is not specified")
	return "", err
}

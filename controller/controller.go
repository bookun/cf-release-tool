package controller

import (
	"io/ioutil"

	"github.com/bookun/cf-release-tool/entity"
	"github.com/bookun/cf-release-tool/usecase"
	yaml "gopkg.in/yaml.v2"
)

type CurrentInfoGetter interface {
	GetCurrentSpace() (string, error)
	GetCurrentOrg() (string, error)
}

type Controller struct {
	InputPort    usecase.InputPort
	ManifestFile string
	Branch       string
}

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

func (c *Controller) Release() error {
	m, err := c.getManifest()
	if err != nil {
		return err
	}
	targetApp := m.Applications[0]
	domain := targetApp.Env.Domain
	host := targetApp.Env.Host
	entity := entity.Deploy{
		Org:          targetApp.Env.Org,
		Space:        targetApp.Env.Space,
		App:          targetApp.Name,
		ManifestFile: c.ManifestFile,
		MaterialDir:  targetApp.Env.Material,
		Branch:       c.Branch,
	}
	if err := c.InputPort.BlueGreenDeployment(entity, domain, host); err != nil {
		return err
	}
	return nil
}

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

package entity

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Manifest has information of manifest.yml
type Manifest struct {
	Applications []App `yaml:"applications"`
}

type App struct {
	Name       string `yaml:"name"`
	Instance   int    `yaml:"instances"`
	Memory     string `yaml:"momery"`
	Buildpack  string `yaml:"buildpack"`
	NoHostName bool   `yaml:"no-hostname"`
	NoRoute    bool   `yaml:"no-route"`
	Host       string `yaml:"host"`
	Domain     string `yaml:"domain"`
	Env        struct {
		Org      string            `yaml:"ORG"`
		Space    string            `yaml:"SPACE"`
		Branch   string            `yaml:"BRANCH"`
		TimeZone string            `yaml:"TZ"`
		Lang     string            `yaml:"LANG"`
		Copy     map[string]string `yaml:"COPY"`
		TestUp   map[string]string `yaml:"TESTUP"`
	} `yaml:"env"`
}

// NewManifest
func NewManifest(manifestFile, branch, host string) (*Manifest, error) {
	data, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return nil, err
	}
	m := &Manifest{}
	if err := yaml.Unmarshal(data, m); err != nil {
		return nil, err
	}
	for _, app := range m.Applications {
		if app.Env.TestUp == nil {
			continue
		}
		parameters := []string{"domain", "host"}
		for _, parameter := range parameters {
			if _, ok := app.Env.TestUp[parameter]; !ok {
				err := fmt.Errorf("the %s for testup is not be set", parameter)
				return nil, err
			}
		}
	}
	if branch != "" {
		for _, app := range m.Applications {
			app.Env.Branch = branch
		}
	}

	if host != "" {
		for _, app := range m.Applications {
			app.Host = host
		}
	}
	return m, nil
}

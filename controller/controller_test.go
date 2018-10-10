package controller

import (
	"errors"
	"testing"

	"github.com/bookun/cf-release-tool/entity"
)

type MockInputPort struct {
}

func (m *MockInputPort) BlueGreenDeployment(entity entity.Deploy, domain, host string) error {
	return nil
}
func (m *MockInputPort) Deployment(entity entity.Deploy, domain, host string) error {
	return nil
}

type MockInputPort2 struct {
}

func (m *MockInputPort2) BlueGreenDeployment(entity entity.Deploy, domain, host string) error {
	return errors.New("BlueGreenDeployment error")
}
func (m *MockInputPort2) Deployment(entity entity.Deploy, domain, host string) error {
	return errors.New("Deployment error")
}

type MockInfoGetter struct{}

func (m *MockInfoGetter) AppExists(app string) error {
	if app == "nothing" {
		return errors.New("app nothing mock error")
	}
	return nil
}

func TestRelease(t *testing.T) {
	c1 := &Controller{InputPort: &MockInputPort{}, InfoGetter: &MockInfoGetter{}, ManifestFile: "../testdata/manifest1.yml", Branch: "master"}
	c2 := &Controller{InputPort: &MockInputPort{}, InfoGetter: &MockInfoGetter{}, ManifestFile: "../testdata/manifest2.yml", Branch: "master"}
	c3 := &Controller{InputPort: &MockInputPort2{}, InfoGetter: &MockInfoGetter{}, ManifestFile: "../testdata/manifest1.yml", Branch: "master"}
	c4 := &Controller{InputPort: &MockInputPort2{}, InfoGetter: &MockInfoGetter{}, ManifestFile: "../testdata/manifest2.yml", Branch: "master"}
	cases := []struct {
		name       string
		controller *Controller
		expect     error
	}{
		{"blue-gleen deployment success", c1, nil},
		{"deployt success", c2, nil},
		{"blue-gleen deployment failed", c3, errors.New("BlueGreenDeployment error")},
		{"deployment failed", c4, errors.New("Deployment error")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if actual := c.controller.Release(); c.expect != actual {
				if c.expect != nil {
					if c.expect.Error() != actual.Error() {
						t.Errorf(
							"want Release() = %v, got %v",
							c.expect, actual)
					}
				}
			}
		})
	}
	// TODO: check parse error
}

package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/bookun/cf-release-tool/client"
	"github.com/bookun/cf-release-tool/entity"
	"github.com/bookun/cf-release-tool/manager"
	"github.com/k0kubun/pp"
)

const (
	host   = "test_host"
	domain = "test_domain"
)

type MockManager struct{}

func (m *MockManager) Init(materialDir, branch, org, space string) error {
	return errors.New("Init error")
}
func (m *MockManager) GreenPush(app, manifestFile, domain, host string) (string, error) {
	return "", errors.New("Green push error")
}
func (m *MockManager) Push(app, manifestFile, domain, host string) error {
	return errors.New("Push error")
}
func (m *MockManager) Exchange(app, blueApp string) (string, error) {
	return "", errors.New("Exchange error")
}

func (m *MockManager) BlueDelete(app, domain, host string) error {
	return errors.New("BlueDelete error")
}

func TestBlueGreenDeployment(t *testing.T) {
	entity := entity.Deploy{
		Org:          "test_org",
		Space:        "test_space",
		App:          "test_app",
		ManifestFile: "test_manifestfile",
		MaterialDir:  "test_materialdir",
		Branch:       "test_branch",
	}
	buf := new(bytes.Buffer)
	bufExpected := new(bytes.Buffer)
	writeBlueGreenDeployment(bufExpected, entity)
	client := client.NewDummyClient(buf)
	manager := manager.NewManager(client)
	usecase1 := NewUsecase(manager)
	usecase2 := NewUsecase(&MockManager{})
	cases := []struct {
		name      string
		usecase   *Usecase
		expectErr error
	}{
		{"success", usecase1, nil},
		{"failed", usecase2, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if actual := c.usecase.BlueGreenDeployment(entity, domain, host); c.expectErr != actual {
				if c.expectErr != nil {
					if c.expectErr.Error() != actual.Error() {
						t.Errorf(
							"want BlueGreenDeployment() = %v, got %v",
							c.expectErr, actual)
					}
				} else {
					if buf.String() != bufExpected.String() {
						pp.Println(bufExpected.String())
						pp.Println(buf.String())
						t.Errorf(
							"want BlueGreenDeployment() = %v, got %v",
							bufExpected.String(), buf.String())
					}
				}
			}
		})
	}
	// TODO: check parse error

}

func writeBlueGreenDeployment(writer io.Writer, entity entity.Deploy) {
	fmt.Fprintf(writer, "rm -fr ./.bp-config\n")
	fmt.Fprintf(writer, "cp -r %s ./.bp-config\n", entity.MaterialDir)
	fmt.Fprintf(writer, "git checkout %s\n", entity.Branch)
	fmt.Fprintf(writer, "git pull origin %s\n", entity.Branch)
	fmt.Fprintf(writer, "target -o %s -s %s\n", entity.Org, entity.Space)
	fmt.Fprintf(writer, "push %s -f %s\n", entity.App+"_green", entity.ManifestFile)
	fmt.Fprintf(writer, "map-route %s %s --hostname %s\n", entity.App+"_green", domain, host)
	fmt.Fprintf(writer, "rename %s to %s\n", entity.App, entity.App+"_blue")
	fmt.Fprintf(writer, "rename %s to %s\n", entity.App+"_green", entity.App)
	fmt.Fprintf(writer, "unmap-route %s %s --hostname %s\n", entity.App+"_blue", domain, host)
	fmt.Fprintf(writer, "delete %s\n", entity.App+"_blue")
}

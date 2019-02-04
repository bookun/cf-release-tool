package manager

import (
	"errors"
)

// Client has aggregation of methods.
// These methods are implemented in client package.
type Client interface {
	Init(copyTargets map[string]string, branch, org, space string) error
	Push(app, baseApp string) error
	Rename(from, to string) error
	Stop(app string) error
	Delete(app string) error
	MapRoute(app, domain, host string) error
	UnMapRoute(app string) error
	TestUp(app, domain string) (bool, error)
	CreateBlueName(app string) (string, error)
}

// Manager has client.Client
type Manager struct {
	client Client
}

// NewManager init Manager
func NewManager(client Client) *Manager {
	return &Manager{
		client: client,
	}
}

// Init call client.Init
func (m *Manager) Init(copyTargets map[string]string, branch, org, space string) error {
	if err := m.client.Init(copyTargets, branch, org, space); err != nil {
		return err
	}
	return nil
}

// GreenPush push newApp and map-route to the app.
// This method is called if there is app you want to deploy in your cloudfoundry space.
func (m *Manager) GreenPush(app, manifestFile, domain, host string) (string, error) {
	newApp := app + "-green"
	if err := m.client.Push(newApp, manifestFile); err != nil {
		return "", err
	}
	confirm, err := m.client.TestUp(newApp, domain)
	if err != nil {
		return "", err
	}
	if confirm {
		if err := m.client.MapRoute(newApp, domain, host); err != nil {
			return "", err
		}
		return newApp, nil
	}
	err = errors.New("cancel release")
	return "", err
}

// Push push newApp and map-route to the app.
// This method is called if there is "not" app you want to deploy in your cloudfoundry space.
func (m *Manager) Push(app, manifestFile, domain, host string) error {
	if err := m.client.Push(app, manifestFile); err != nil {
		return err
	}
	if err := m.client.MapRoute(app, domain, host); err != nil {
		return err
	}
	return nil
}

// Exchange exchange name between app andd blueApp.
func (m *Manager) Exchange(app, blueApp string) (string, error) {
	//oldApp := app + "_blue"
	oldApp, err := m.client.CreateBlueName(app)
	if err != nil {
		return "", err
	}
	if err := m.client.Rename(app, oldApp); err != nil {
		return "", err
	}
	if err := m.client.Rename(blueApp, app); err != nil {
		return "", err
	}
	return oldApp, nil
}

// BlueDelete delete old app.
func (m *Manager) BlueDelete(app, domain, host string) error {
	if err := m.client.UnMapRoute(app); err != nil {
		return err
	}
	if err := m.client.Stop(app); err != nil {
		return err
	}
	// TODO: 本当はここで商用にエラーが多発してないかとかチェックしたい
	//appType := strings.Split(app, "-")[0]
	if err := m.client.Delete(app); err != nil {
		return err
	}
	return nil
}

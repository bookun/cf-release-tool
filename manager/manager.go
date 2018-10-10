package manager

// Clientは clientパッケージで実装される
type Client interface {
	Init(materialDir, branch, org, space string) error
	Push(app, baseApp string) error
	Rename(from, to string) error
	Delete(app string) error
	MapRoute(app, domain, host string) error
	UnMapRoute(app, domain, host string) error
}

type Manager struct {
	client Client
}

func NewManager(client Client) *Manager {
	return &Manager{
		client: client,
	}
}

func (m *Manager) Init(materialDir, branch, org, space string) error {
	if err := m.client.Init(materialDir, branch, org, space); err != nil {
		return err
	}
	return nil
}

func (m *Manager) BluePush(app, manifestFile, domain, host string) (string, error) {
	newApp := app + "_blue"
	if err := m.client.Push(newApp, manifestFile); err != nil {
		return "", err
	}
	if err := m.client.MapRoute(newApp, domain, host); err != nil {
		return "", err
	}
	return newApp, nil

}

func (m *Manager) Exchange(app, blueApp string) (string, error) {
	oldApp := app + "_green"
	if err := m.client.Rename(app, oldApp); err != nil {
		return "", err
	}
	if err := m.client.Rename(blueApp, app); err != nil {
		return "", err
	}
	return oldApp, nil
}

func (m *Manager) GreenDelete(app, domain, host string) error {
	if err := m.client.UnMapRoute(app, domain, host); err != nil {
		return err
	}
	// TODO: 本当はここで商用にエラーが多発してないかとかチェックしたい

	if err := m.client.Delete(app); err != nil {
		return err
	}
	return nil
}

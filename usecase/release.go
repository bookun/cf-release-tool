package usecase

import "github.com/bookun/cf-release-tool/entity"

// CfManager is aggregation of methods.
// Each methods is implemented in manager package.
type CfManager interface {
	Init(materialDir, branch, org, space string) error
	GreenPush(app, manifestFile, domain, host string) (string, error)
	Push(app, manifestFile, domain, host string) error
	Exchange(app, blueApp string) (string, error)
	BlueDelete(app, domain, host string) error
}

// InputPort is called by controller
type InputPort interface {
	BlueGreenDeployment(entity entity.Deploy, domain, host string) error
	Deployment(entity entity.Deploy, domain, host string) error
}

// Usecase has CFManeger
type Usecase struct {
	client CfManager
}

// NewUsecase init Usecase.
func NewUsecase(manager CfManager) *Usecase {
	return &Usecase{
		client: manager,
	}
}

// BlueGreenDeployment executes BlueGreenDeployment.
// At first, it deploys new app(green app) and map-route.
// Then, it exchanges name between green app and one that already deployed (blue app).
// Finally, it deletes blue app.
func (u *Usecase) BlueGreenDeployment(entity entity.Deploy, domain, host string) error {
	if err := u.client.Init(entity.MaterialDir, entity.Branch, entity.Org, entity.Space); err != nil {
		return err
	}
	greenApp, err := u.client.GreenPush(entity.App, entity.ManifestFile, domain, host)
	if err != nil {
		return err
	}
	blueApp, err := u.client.Exchange(entity.App, greenApp)
	if err != nil {
		return err
	}

	if err := u.client.BlueDelete(blueApp, domain, host); err != nil {
		return err
	}
	return nil
}

// Deployment executes deploy an app.
func (u *Usecase) Deployment(entity entity.Deploy, domain, host string) error {
	if err := u.client.Init(entity.MaterialDir, entity.Branch, entity.Org, entity.Space); err != nil {
		return err
	}
	if err := u.client.Push(entity.App, entity.ManifestFile, domain, host); err != nil {
		return err
	}
	return nil
}

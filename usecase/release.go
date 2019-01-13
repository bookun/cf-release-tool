package usecase

// CfManager is aggregation of methods.
// Each methods is implemented in manager package.
type CfClient interface {
	Init() error
	Push(name string) error
	Rename(from, to string) error
	//RenameFrom(from string) error
	//RenameTo(to string) error
	Stop(app string) error
	DeleteApps(app string) error
	MapRoute(app string) error
	UnMapRoute(app string) error
	//TestUp(app, domain string) (bool, error)
	CreateBlueName(app string) (string, error)
}

// InputPort is called by controller
type InputPort interface {
	BlueGreenDeployment() error
	Deployment() error
}

// Usecase has CFManeger
type Usecase struct {
	appName string
	client CfClient
}

// NewUsecase init Usecase.
func NewUsecase(appName string, client CfClient) *Usecase {
	return &Usecase{
		appName:appName,
		client: client,
	}
}

// BlueGreenDeployment executes BlueGreenDeployment.
// At first, it deploys new app(green app) and map-route.
// Then, it exchanges name between green app and one that already deployed (blue app).
// Finally, it deletes blue app.
func (u *Usecase) BlueGreenDeployment() error {
	if err := u.client.Init(); err != nil {
		return err
	}
	greenAppName := "green-"+u.appName
	blueAppName, err := u.client.CreateBlueName(u.appName)
	if err != nil {
		return err
	}
	if err := u.client.Push(greenAppName); err != nil {
		return err
	}
	if err := u.client.MapRoute(greenAppName); err != nil {
		return err
	}
	if err := u.client.UnMapRoute(greenAppName); err != nil {
		return err
	}
	if err := u.client.Rename(u.appName, blueAppName); err != nil {
		return err
	}
	if err := u.client.Rename(greenAppName, u.appName); err != nil {
		return err
	}
	if err := u.client.MapRoute(u.appName); err != nil {
		return err
	}
	if err := u.client.UnMapRoute(blueAppName); err != nil {
		return err
	}
	if err := u.client.Stop(blueAppName); err != nil {
		return err
	}
	if err := u.client.DeleteApps(u.appName); err != nil {
		return err
	}
	return nil
}

// Deployment executes deploy an app.
func (u *Usecase) Deployment() error {
	if err := u.client.Init(); err != nil {
		return err
	}
	if err := u.client.Push(u.appName); err != nil {
		return err
	}
	if err := u.client.MapRoute(u.appName); err != nil {
		return err
	}
	return nil
}

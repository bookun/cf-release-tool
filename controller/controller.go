package controller

import (
	"github.com/bookun/cf-release-tool/usecase"
)

// CurrentInfoGetter interface has AppExists method.
// AppExists method is implemented in client package.
type CurrentInfoGetter interface {
	AppExists() error
}

// Controller decides usecase, BlueGreenDeployment or normal Deployment.
type Controller struct {
	InputPort    usecase.InputPort
	InfoGetter   CurrentInfoGetter
}


// Release executes deployment.
// If there is application you want to release, this method executes normal deployment.
// Else, it executes BlueGreenDeployment.
func (c *Controller) Release() error {
	if c.InfoGetter.AppExists() != nil {
		if err := c.InputPort.Deployment(); err != nil {
			return err
		}
	} else {
		if err := c.InputPort.BlueGreenDeployment(); err != nil {
			return err
		}
	}
	return nil
}

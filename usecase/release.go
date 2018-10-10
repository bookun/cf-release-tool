package usecase

import "github.com/bookun/cf-release-tool/entity"

// CfManager の関数は managerパッケージで実装される
type CfManager interface {
	Init(materialDir, branch, org, space string) error
	BluePush(app, manifestFile, domain, host string) (string, error)
	Exchange(app, blueApp string) (string, error)
	GreenDelete(app, domain, host string) error
}

// InputPort defines inputPort
type InputPort interface {
	BlueGreenDeployment(entity entity.Deploy, domain, host string) error
}

// Usecase は CfManagerのもつメソッドを組み立ててユースケースを実行する
type Usecase struct {
	client CfManager
}

func NewUsecase(manager CfManager) *Usecase {
	return &Usecase{
		client: manager,
	}
}

// BlueGreenDeployment は BlueGreenDeploymentをする
func (u *Usecase) BlueGreenDeployment(entity entity.Deploy, domain, host string) error {
	// TODO: 最初の色がどっちか忘れた
	if err := u.client.Init(entity.MaterialDir, entity.Branch, entity.Org, entity.Space); err != nil {
		return err
	}
	blueApp, err := u.client.BluePush(entity.App, entity.ManifestFile, domain, host)
	if err != nil {
		return err
	}
	greenApp, err := u.client.Exchange(entity.App, blueApp)
	if err != nil {
		return err
	}

	if err := u.client.GreenDelete(greenApp, domain, host); err != nil {
		return err
	}
	return nil
}

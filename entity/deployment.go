package entity

// Deploy is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type Deployment struct {
	Org    string
	Space  string
	Branch string
}

func NewDeployment(org, space, branch string) *Deployment {
	return &Deployment{
		Org:    org,
		Space:  space,
		Branch: branch,
	}
}

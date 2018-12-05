package entity

// Deploy is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type Deploy struct {
	Org          string
	Space        string
	App          string
	ManifestFile string
	MaterialDir  string
	Branch       string
	EnvFile      string
}

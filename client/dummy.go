package client

import (
	"errors"
	"fmt"
	"io"
)

// DummyClient is mock about client.Client
type DummyClient struct {
	Output io.Writer
}

// NewDummyClient is mock about NewClient in client.go
func NewDummyClient(output io.Writer) *DummyClient {
	return &DummyClient{
		Output: output,
	}
}

// Init is mock about Init in client.go
func (c *DummyClient) Init(materialDir, branch, org, space string) error {
	fmt.Fprintf(c.Output, "rm -fr ./.bp-config\n")
	fmt.Fprintf(c.Output, "cp -r %s ./.bp-config\n", materialDir)
	fmt.Fprintf(c.Output, "git checkout %s\n", branch)
	fmt.Fprintf(c.Output, "git pull origin %s\n", branch)
	fmt.Fprintf(c.Output, "target -o %s -s %s\n", org, space)
	return nil
}

// Push is mock about Push in client.go
func (c *DummyClient) Push(app, manifestFile string) error {
	fmt.Fprintf(c.Output, "push %s -f %s\n", app, manifestFile)
	return nil
}

// Rename is mock about Rename in client.go
func (c *DummyClient) Rename(oldApp, newApp string) error {
	fmt.Fprintf(c.Output, "rename %s to %s\n", oldApp, newApp)
	return nil
}

// Delete is mock about Delete in client.go
func (c *DummyClient) Delete(app string) error {
	fmt.Fprintf(c.Output, "delete %s\n", app)
	return nil
}

// MapRoute is mock about MapRoute in client.go
func (c *DummyClient) MapRoute(app, domain, host string) error {
	if host != "" {
		fmt.Fprintf(c.Output, "map-route %s %s --hostname %s\n", app, domain, host)
	} else {
		fmt.Fprintf(c.Output, "map-route %s %s\n", app, domain)
	}
	return nil
}

// UnMapRoute is mock about UnMapRoute in client.go
func (c *DummyClient) UnMapRoute(app, domain, host string) error {
	if host != "" {
		fmt.Fprintf(c.Output, "unmap-route %s %s --hostname %s\n", app, domain, host)
	} else {
		fmt.Fprintf(c.Output, "unmap-route %s %s\n", app, domain)
	}
	return nil
}

// AppExists is mock about AppExists in client.go
func (c *DummyClient) AppExists(app string) error {
	if app == "nothing" {
		return errors.New("nothing")
	}
	return nil
}

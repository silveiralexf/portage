package golang

import "dagger.io/dagger"

func SetupTparse(c *dagger.Container) *dagger.Container {
	goInstall := []string{"go", "install", "-v", "github.com/mfridman/tparse@latest"}

	return c.WithExec(goInstall)
}

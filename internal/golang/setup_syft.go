package golang

import "dagger.io/dagger"

func SetupSyft(c *dagger.Container) *dagger.Container {
	goInstall := []string{"go", "install", "-v", "github.com/anchore/syft/cmd/syft@latest"}

	return c.WithExec(goInstall)
}

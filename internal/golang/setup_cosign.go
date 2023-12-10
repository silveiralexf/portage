package golang

import "dagger.io/dagger"

func SetupCosign(c *dagger.Container) *dagger.Container {
	goInstall := []string{"go", "install", "-v", "github.com/sigstore/cosign/v2/cmd/cosign@latest"}

	return c.WithExec(goInstall)
}

package golang

import (
	"fmt"

	"dagger.io/dagger"
)

func SetupGo(c *dagger.Container) *dagger.Container {
	version := "1.21.4"
	arch := "$(dpkg --print-architecture)"
	curl := fmt.Sprintf("curl -o go_linux.tar.gz -L https://go.dev/dl/go%s.linux-%s.tar.gz", version, arch)
	return c.WithExec([]string{"sh", "-c", curl}).
		WithExec([]string{"tar", "-C", "/usr/local", "-xvf", "go_linux.tar.gz"}).
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "amd64").
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("PATH", "/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
}

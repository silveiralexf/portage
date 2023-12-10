package builder

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/silveiralexf/gomakeme/internal/golang"
)

// echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

func NewUbuntuBuilder(client *dagger.Client, source *dagger.Directory) *dagger.Container {
	return client.Container().Pipeline("build").
		From("ubuntu:jammy@sha256:83f0c2a8d6f266d687d55b5cb1cb2201148eb7ac449e4202d9646b9083f1cee0").
		WithEnvVariable("DOCKER_DEFAULT_PLATFORM", "linux/amd64").
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "amd64").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "gpg", "nix-bin", "upx-ucl"}).
		With(NewDockerContainer).
		WithExec([]string{"adduser", "-q", "nonroot"}).
		WithExec([]string{"mkdir", "/nix"}).
		WithExec([]string{"chown", "nonroot", "/nix"}).
		WithMountedDirectory("/src", source).
		With(golang.SetupGo).WithUser("nonroot").WithWorkdir("/src").
		WithNewFile("/usr/local/sbin/docker", dagger.ContainerWithNewFileOpts{
			Contents: `#!/bin/sh
		DOCKER_HOST=tcp://localhost:2375 /usr/bin/docker $@`,
			Permissions: 0o777,
		}).
		With(golang.SetupCosign).
		WithServiceBinding("localhost", WithDinD(client)).
		WithEnvVariable("DOCKER_HOST", "tcp://localhost:2375")
}

// WithMountedCache("/gomods", client.CacheVolume("gomodcache")).
// WithExec([]string{"chown", "nonroot", "-R", "/gomods"}).
// WithEnvVariable("GOMODCACHE", "/gomods").
// WithMountedCache("/tmp/gocache", client.CacheVolume("gobuildcache")).
// WithExec([]string{"chown", "nonroot", "-R", "/tmp/gocache"}).
// WithEnvVariable("GOCACHE", "/tmp/gocache").
// WithExec([]string{"chown", "-R", "nonroot", "/gomods"}).

func WithDinD(client *dagger.Client) *dagger.Service {
	container := client.Container().Pipeline("dockerd").From("docker:20-dind").WithExposedPort(2375).WithEnvVariable(
		"DOCKER_DEFAULT_PLATFORM", "linux/amd64",
	).WithExec([]string{
		"dockerd", "--log-level=trace", "--host=tcp://0.0.0.0:2375", "--tls=false",
	},
		dagger.ContainerWithExecOpts{InsecureRootCapabilities: true},
	)

	args := []dagger.BuildArg{}
	args = append(args, dagger.BuildArg{Name: "", Value: ""})
	container.Build(&dagger.Directory{}, dagger.ContainerBuildOpts{BuildArgs: args})
	return container.AsService()
}

func NewDockerContainer(c *dagger.Container) *dagger.Container {
	version := "20.10.24"
	version_apt := fmt.Sprintf("5:%s~3-0~ubuntu-jammy", version)
	install := fmt.Sprintf("apt-get install -y docker-ce=%s docker-ce-cli=%s containerd.io docker-buildx-plugin docker-compose-plugin", version_apt, version_apt)
	return c.
		WithExec([]string{"apt-get", "install", "ca-certificates", "curl", "gnupg"}).
		WithExec([]string{"install", "-m", "0755", "-d", "/etc/apt/keyrings"}).
		WithExec([]string{"sh", "-c", "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg"}).
		WithExec([]string{"chmod", "a+r", "/etc/apt/keyrings/docker.gpg"}).
		WithExec([]string{"sh", "-c", `echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null`}).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"sh", "-c", install})
}

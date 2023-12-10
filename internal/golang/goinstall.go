package golang

import (
	"fmt"

	"dagger.io/dagger"
)

// GoInstall will go install a given package given a set of attribute flags, package name and version
func GoInstall(c *dagger.Container, packageName, packageVersion string, extraFlags []string) *dagger.Container {
	goInstall := []string{}
	goInstall = append(goInstall, "go", "install", "-v")
	goInstall = append(goInstall, extraFlags...)
	goInstall = append(goInstall, fmt.Sprintf("%v@%v", packageName, packageVersion))
	return c.WithExec(goInstall)
}

package builder

import (
	"fmt"

	"dagger.io/dagger"
)

type BaseBuilder struct {
	*dagger.Client
	*dagger.Directory
}

func (b *BaseBuilder) New(opts BuilderOptions) error {
	b.Client = &dagger.Client{}
	b.Client.Container().Pipeline(opts.targetName).From(opts.FromUbuntuJammy().baseOSImage)

	return nil
}

type BuilderOptions struct {
	baseOSImage    string
	targetTechType string
	targetName     string
}

// FromUbuntuJammy defines FROM directive as the Base OS source to use Ubuntu Jammy (22.04)
func (opts *BuilderOptions) FromUbuntuJammy() *BuilderOptions {
	opts.baseOSImage = "ubuntu:jammy@sha256:83f0c2a8d6f266d687d55b5cb1cb2201148eb7ac449e4202d9646b9083f1cee0"
	return opts
}

func (opts *BuilderOptions) WithGolang() *BuilderOptions {
	opts.targetTechType = "golang"
	return opts
}

func (opts *BuilderOptions) WithTarget(targetName string) *BuilderOptions {
	opts.targetName = targetName
	return opts
}

func NewMain() {
	opts := BuilderOptions{}
	opts.FromUbuntuJammy().WithGolang().WithTarget("build")

	b := BaseBuilder{}
	err := b.New(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(opts)
	// Initialize a Count object and pass it to WriteLog().

}

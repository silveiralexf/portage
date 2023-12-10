/*
Copyright © 2023 silveiralexf

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/spf13/cobra"

	"github.com/silveiralexf/gomakeme/internal/builder"
)

// golangCmd represents the golang command
var golangCmd = &cobra.Command{
	Use:   "golang",
	Short: "Builds a golang project",
	Run: func(cmd *cobra.Command, args []string) {
		err := doGolang()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func doGolang() error {
	ctx := context.Background()

	fmt.Println("Starting client...")
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Directory(".",
		dagger.HostDirectoryOpts{
			Exclude: []string{
				".git/",
			},
		})

	fmt.Println("Running CI steps...")
	builder, err := builder.NewUbuntuBuilder(client, src).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create CI environment: %v", err)
	}

	builder = builder.WithMountedDirectory("/src", builder.Directory("/src")).WithUser("root").
		WithEnvVariable("DOCKER_DEFAULT_PLATFORM", "linux/amd64").
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "amd64").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build", "-o", "gomakeme"}).
		WithExec([]string{"sh", "-c", "go test -failfast -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./... -run . -timeout=15m"})

	out, err := builder.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed during CI pipeline stages: %v", err)
	}
	fmt.Println(out)

	// Export built binary
	_, err = builder.File("./gomakeme").Export(ctx, "./gomakeme")
	if err != nil {
		return fmt.Errorf("failure during CI pipeline step: %+v", err)
	}

	ref, err := builder.Publish(context.Background(), "silveiralexf/gomakeme:latest", dagger.ContainerPublishOpts{})
	if err != nil {
		return err
	}
	fmt.Println("REF", ref)
	return nil
}

func init() {
	rootCmd.AddCommand(golangCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// golangCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// golangCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

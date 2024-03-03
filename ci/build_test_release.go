package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// use a node:16-slim container
	// mount the source code directory on the host
	// at /src in the container
	source := client.Container().
		From("golang:1.21").
		WithDirectory("/src", client.Host().Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"ci/"},
		}))

	geese := []string{"darwin", "linux", "windows"}
	goarch := "amd64"

	// set the working directory in the container
	// install application dependencies
	runner := source.WithWorkdir("/src").
		WithExec([]string{"go", "mod", "tidy"})

	runner = source.WithWorkdir("/src/src_code/go_src")
	// run application tests

	test := runner.WithExec([]string{"go", "test"})
	
	buildDir := test.Directory("/src/dist")

	for _, goos := range geese {
		path := fmt.Sprintf("/dist/")
		filename := fmt.Sprintf("/dist/cocommit_go-%s", goos) 
		// build application
		// write the build output to the host
		build := test.
			WithEnvVariable("GOOS", goos).
			WithEnvVariable("GOARCH", goarch).
			WithExec([]string{"go", "build", "-o", filename})

		buildDir = buildDir.WithDirectory(path, build.Directory(path))

	}

	_, err = buildDir.Export(ctx, ".")
	if err != nil {
		panic(err)
	}
	e, err := buildDir.Entries(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("build dir contents:\n %s\n", e)
}
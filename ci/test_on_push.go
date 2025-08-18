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

	goCache := client.CacheVolume("golang")

	// use a node:16-slim container
	// mount the source code directory on the host
	// at /src in the container
	source := client.Container().
		From("golang:1.24").
		WithDirectory("/src_d", client.Host().Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"build/"},
		})).WithMountedCache("/src_d/dagger_dep_cache/go_dep", goCache)

		// set the working directory in the container
		// install application dependencies
	runner := source.WithWorkdir("/src_d/src").
		WithExec([]string{"go", "mod", "tidy"}).WithEnvVariable("CI", "true")

		// run application tests
	out, err := runner.WithWorkdir("/src_d/src").WithExec([]string{"go", "test", "./..."}).
		Stderr(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

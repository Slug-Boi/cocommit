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
		From("golang:1.23").
		WithDirectory("/src", client.Host().Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"build/"},
		})).WithMountedCache("/src/dagger_dep_cache/go_dep", goCache)

		// set the working directory in the container
		// install application dependencies
	runner := source.WithWorkdir("/src/src_code/go_src/").
		WithExec([]string{"go", "mod", "tidy"})

		// run application tests
	out, err := runner.WithWorkdir("/src/src_code/go_src").WithExec([]string{"go", "test"}).
		Stderr(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
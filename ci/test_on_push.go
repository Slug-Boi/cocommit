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
		WithDirectory("/src_d", client.Host().Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"build/"},
		})).WithMountedCache("/src_d/dagger_dep_cache/go_dep", goCache)

		// set the working directory in the container
		// install application dependencies
	runner := source.WithWorkdir("/src_d/src").
		WithExec([]string{"go", "mod", "tidy"}).WithEnvVariable("CI", "true")

		// run application tests
	out, err := runner.WithWorkdir("/src_d/src").WithExec([]string{"go", "test", "./cmd/utils", "./cmd/tui", "-coverprofile=cover.out"}).
		Stderr(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	out, err = runner.WithExec([]string{"find", "/", "-name", "cover.out"}).Stdout(ctx)
	fmt.Println(out)

	// export the coverage report
	_, err = runner.WithWorkdir("/").File("/src_d/src/cover.out").Export(ctx, "./cover.out")
if err != nil {
    panic(err)
}
}
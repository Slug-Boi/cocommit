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
			Exclude: []string{},
		})).WithMountedCache("/src_d/dagger_dep_cache/go_dep", goCache)

	geese := []string{"darwin", "linux", "windows"}
	goarch := "amd64"

	// set the working directory in the container
	// install application dependencies
	runner := source.WithWorkdir("/src_d/src/").
		WithExec([]string{"go", "mod", "tidy"}).WithEnvVariable("CI", "true")

	// run application tests
	test := runner.WithWorkdir("/src_d/src/").WithExec([]string{"go", "test", "./..."}).WithEnvVariable("CI", "true")

	buildDir := test.Directory("/src_d/src/")

	Coco_var := os.Getenv("Coco_Version")


	for _, goos := range geese {
		path := fmt.Sprintf("/dist/")
		filename := fmt.Sprintf("/dist/cocommit-%s", goos)
		// build application
		// write the build output to the host
		build := test.
			WithEnvVariable("GOOS", goos).
			WithEnvVariable("GOARCH", goarch).
			WithExec([]string{"go", "build", "-o", filename, "-ldflags", "-X github.com/Slug-Boi/cocommit/src/cmd.Coco_Version="+Coco_var}).WithEnvVariable("CI", "true")

		buildDir = buildDir.WithDirectory(path, build.Directory(path))

	}

	// extra step to build for aarch on darwin
	path := fmt.Sprintf("/dist/")
	filename := fmt.Sprintf("/dist/cocommit-darwin-aarch64")

	build := test.
		WithEnvVariable("GOOS", "darwin").
		WithEnvVariable("GOARCH", "arm64").
		WithExec([]string{"go", "build", "-o", filename, "-ldflags", "-X github.com/Slug-Boi/cocommit/src/cmd.Coco_Version="+Coco_var}).WithEnvVariable("CI", "true")

	buildDir = buildDir.WithDirectory(path, build.Directory(path))

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

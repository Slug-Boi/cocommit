Coco_Version ?= "undefined"

build:
	go build -C src -o cocommit -ldflags "-X github.com/Slug-Boi/cocommit/src/cmd.Coco_Version=${Coco_Version}"

build:
	go build -C src -o cocommit 

build-nix:
	GOWORK=off go build -C src -o cocommit

build:
	go build -C src_code/go_src/ -o cocommit 

build-nix:
	GOWORK=off go build -C src_code/go_src/ -o cocommit
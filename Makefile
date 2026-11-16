# A Makefile for the saol.se-cli project

lint:
	echo "TODO: add linting"

test:
	go test ./...

build:
	go build -o bin/saol-cli main.go

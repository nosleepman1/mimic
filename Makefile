# Mimic Makefile

BINARY_NAME=mimic

.PHONY: all build clean test run linux-build

all: build

deps:
	go mod download
	go mod tidy

build:
	go build -ldflags="-s -w" -o bin/$(BINARY_NAME).exe main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)-linux-amd64 main.go

build-darwin:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)-darwin-arm64 main.go

test:
	go test -v -race ./...

clean:
	rm -rf bin/

#!make

# Usage:
# make help

## VARIABLES
BINARY_NAME = gor-argocd
RELEASE_VERSION ?= v0.0.1
# Colour Outputs
GREEN := \033[0;32m
CLEAR := \033[00m

## Targets
build: clean
	@echo "$(GREEN)INFO: Building $(BINARY_NAME) $(CLEAR)"
	go get -v
	go mod tidy
	GOARCH=amd64 GOOS=darwin  go build -o $(BINARY_NAME)-darwin  main.go
	GOARCH=amd64 GOOS=linux   go build -o $(BINARY_NAME)-linux   main.go

run:
	@echo "$(GREEN)RUN: Test run $(CLEAR)"
	./${BINARY_NAME}-darwin || true

build_and_run: build run

release: clean build run
	@echo "$(GREEN)RELEASE: Creating a GitHub release $(CLEAR)"
	gh release create $(RELEASE_VERSION) --generate-notes || true
	gh release upload v0.0.1 $(BINARY_NAME)-linux --clobber
	$(MAKE) clean

release_undo:
	@echo "$(GREEN)UNDO RELEASE: Deleting a GitHub release and the corresponding tag $(CLEAR)"
	gh release delete $(RELEASE_VERSION) --yes || true
	git push --delete origin $(RELEASE_VERSION) || true
	git tag --delete $(RELEASE_VERSION) || true

clean:
	@echo "$(GREEN)CLEAN: Removing all binaries $(CLEAR)"
	go clean
	rm -f $(BINARY_NAME)-darwin $(BINARY_NAME)-linux

help:
	@echo "$(GREEN)HELP: make <command> $(CLEAR)"
	@echo "  build           Build packages for different architectures"
	@echo "  run             Run a command to make sure package is working fine"
	@echo "  build_and_run   Build and Run"
	@echo "  clean           Remove the built packages (if any)"
	@echo "  release         Create a GitHub release"
	@echo "  release_undo    Delete a GitHub release"

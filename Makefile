SHELL := /bin/bash

.PHONY: all check format vet lint build test generate tidy

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"
	@echo "  generate            to generate code"
	@echo "  test                to run test"
	@echo "  integration_test    to run integration test"

# mockgen: go get github.com/golang/mock/mockgen
# golint: go get -u golang.org/x/lint/golint
# definitions: go get -u github.com/aos-dev/go-storage/cmd/definitions
tools := mockgen golint definitions

$(tools):
	@command -v $@ >/dev/null 2>&1 || echo "$@ is not found, please install it."

check: vet lint

format:
	@echo "go fmt"
	@go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

lint: golint
	@echo "golint"
	@golint ./...
	@echo "ok"

generate: definitions mockgen
	@echo "generate code"
	@go generate ./...
	@go fmt ./...
	@echo "ok"

build: generate tidy check
	@echo "build storage"
	@go build ./...
	@echo "ok"

test:
	@echo "run test"
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
	@go tool cover -html="coverage.txt" -o "coverage.html"
	@echo "ok"

tidy:
	@go mod tidy && go mod verify

clean:
	@echo "clean generated files"
	@find . -type f -name 'generated.go' -delete
	@echo "Done"
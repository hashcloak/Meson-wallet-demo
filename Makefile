GOPATH=$(shell go env GOPATH)

.PHONY: default
default: lint

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	go mod tidy

.PHONY: build
build:
	go build ./cmd/wallet

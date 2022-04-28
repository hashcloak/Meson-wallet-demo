GOPATH=$(shell go env GOPATH)

.PHONY: default
default: lint

.PHONY: lint
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.0
	$(GOPATH)/bin/golangci-lint run --timeout 2m0s -e gosec ./...
	go fmt ./...
	go mod tidy

.PHONY: build
build:
	go build ./cmd/wallet

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SHELL := /bin/bash
COMMIT:=$(shell git rev-list -1 HEAD)
VERSION:=$(COMMIT)
DATE:=$(shell date -uR)

BIN_NAME:=auth0-exporter
GOFLAGS:=-mod=readonly
GO_BUILD:=go build $(GOFLAGS)
KO_DOCKER_REPO:=ghcr.io/tfadeyi/auth0-exporter

export KO_DOCKER_REPO=$KO_DOCKER_REPO
# include files with the `// +build mock` annotation
TEST_TAGS:=-tags mock -coverprofile cover.out

.PHONY: build generate docs test test_all lint build-all-platforms clean install-tools licenses docker-build

install-tools:
	go install github.com/swaggo/swag/cmd/swag@v1.8.7
	go install github.com/matryer/moq@v0.2.7
	go install github.com/google/go-licenses@c781b427440f8ea100841eefdd308e660d26d121
	go install github.com/norwoodj/helm-docs/cmd/helm-docs@v1.11.0
	go install github.com/google/ko@v0.13.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2
	go install github.com/goreleaser/goreleaser@v1.16.2

build:
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME) .

./builds/$(BIN_NAME)-$(GOOS)-$(GOARCH):
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME)-$(GOOS)-$(GOARCH) .

build-all-platforms:
	$(MAKE) GOOS=linux   GOARCH=amd64 ./builds/$(BIN_NAME)-linux-amd64
	$(MAKE) GOOS=darwin  GOARCH=amd64 ./builds/$(BIN_NAME)-darwin-amd64
	$(MAKE) GOOS=windows GOARCH=amd64 ./builds/$(BIN_NAME)-windows-amd64

# generates all documentation: openapi spec, helm-docs
docs:
	cd $(ROOT_DIR) && go generate ./... && \
	helm-docs --chart-search-root=deploy/charts/

# used to generate struct mocks
generate:
	cd $(ROOT_DIR) && go generate ./...

test: build
	cd $(ROOT_DIR) &&  go test $(GOFLAGS) $(TEST_TAGS) ./...

lint:
	golangci-lint run && \
	goreleaser check

test_all: test lint

coverage: test
	go tool cover -html=cover.out -o cover.html
	@echo "open ./cover.html to see coverage"

clean:
	cd $(ROOT_DIR) && \
	rm -rf ./builds && \
	rm -rf kodata

# generates the licenses used by the tool
licenses:
	rm -rf kodata
	go-licenses save . --save_path="kodata/licenses"

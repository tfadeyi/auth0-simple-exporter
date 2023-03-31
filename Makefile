ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SHELL := /bin/bash
COMMIT:=$(shell git rev-list -1 HEAD)
VERSION:=$(COMMIT)
DATE:=$(shell date -uR)

BIN_NAME:=auth0-exporter
GOFLAGS:=-mod=readonly
GO_BUILD:=go build $(GOFLAGS)

# include files with the `// +build mock` annotation
TEST_TAGS:=-tags mock -coverprofile cover.out

.PHONY: build generate test vet build-all-platforms clean install-tools licenses generate docker-build

install-tools:
	go install github.com/swaggo/swag/cmd/swag@v1.8.7
	go install github.com/matryer/moq@v0.2.7
	go install github.com/google/go-licenses@c781b427440f8ea100841eefdd308e660d26d121
	go install github.com/norwoodj/helm-docs/cmd/helm-docs@v1.11.0
	go install github.com/google/ko@latest

build:
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME) .

generate:
	cd $(ROOT_DIR) && go generate ./... && \
	helm-docs --chart-search-root=deploy/charts/

test: build
	cd $(ROOT_DIR) &&  go test $(GOFLAGS) $(TEST_TAGS) ./...

coverage: test
	go tool cover -html=cover.out -o cover.html
	@echo "open ./cover.html to see coverage"

vet:
	cd $(ROOT_DIR) && go vet $(GOFLAGS) ./...

stop:
	docker-compose down --volumes

./builds/$(BIN_NAME)-$(GOOS)-$(GOARCH):
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME)-$(GOOS)-$(GOARCH) .

build-all-platforms:
	$(MAKE) GOOS=linux   GOARCH=amd64 ./builds/$(BIN_NAME)-linux-amd64
	$(MAKE) GOOS=darwin  GOARCH=amd64 ./builds/$(BIN_NAME)-darwin-amd64
	$(MAKE) GOOS=windows GOARCH=amd64 ./builds/$(BIN_NAME)-windows-amd64

docker-build:
	ko build -t latest -t $COMMIT \
	--platform=linux/amd64

clean:
	cd $(ROOT_DIR) && \
	rm -rf ./builds

licenses:
	rm -rf kodata
	go-licenses save . --save_path="kodata/licenses"

swagger:
	go run main.go swagger
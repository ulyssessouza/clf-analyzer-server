# Configuration
#FIXME Consider using environment variables
APP_PORT=8000
GIN_PORT=3000

# Go parameters
GOCMD=go

GOBUILD_ARGS=-o $(BINARY_TARGET_PATH) -v -gcflags='-N -l'
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BUILD_TARGET_PATH=dist
BINARY_NAME=clf-analyzer-server
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME)_win

GIN_BUILD_ARGS="$(GOBUILD_ARGS)"

BINARY_TARGET_PATH=$(BUILD_TARGET_PATH)/$(BINARY_NAME)
BINARY_TARGET_UNIX_PATH=$(BUILD_TARGET_PATH)/$(BINARY_UNIX)
BINARY_TARGET_WINDOWS_PATH=$(BUILD_TARGET_PATH)/$(BINARY_WINDOWS)

ensure-progs: ensure-dep ensure-gin
	@echo ensure-progs

all: test rundev

ensure:
	dep ensure

build: ensure-progs ensure clean goformat
	$(GOBUILD) $(GOBUILD_ARGS)

test: build
	$(GOTEST) -v github.com/ulyssessouza/clf-analyzer-server/core
	$(GOTEST) -v github.com/ulyssessouza/clf-analyzer-server/data
	$(GOTEST) -v github.com/ulyssessouza/clf-analyzer-server/http

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_TARGET_PATH)

run: build
	$(BINARY_TARGET_PATH)

# Cross compilation
build-linux: goformat
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_TARGET_UNIX_PATH) -v

docker-build: check-env-GOPATH goformat
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/ulyssessouza/clf-analyzer-server golang:latest go build -o "$(BINARY_TARGET_UNIX_PATH)"

# Util
check-env-%:
	@ if [ "${${*}}" = "" ]; then \
		@echo "Environment variable $* not set"; \
		exit 1; \
	fi

goformat:
	go fmt .

ensure-dep:
ifeq (, $(shell which dep))
	go get -u github.com/golang/dep/cmd/dep
endif
	@echo ensure dep

ensure-gin:
ifeq (, $(shell which gin))
	go get -u github.com/codegangsta/gin
endif
	@echo ensure gin

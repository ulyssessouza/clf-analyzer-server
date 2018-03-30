# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BUILD_TARGET_PATH=target
BINARY_NAME=clf-anayzer-server
BINARY_UNIX=$(BINARY_NAME)_unix

BINARY_TARGET_PATH=$(BUILD_TARGET_PATH)/$(BINARY_NAME)
BINARY_TARGET_UNIX_PATH=$(BUILD_TARGET_PATH)/$(BINARY_UNIX)
BINARY_TARGET_WINDOWS_PATH=$(BUILD_TARGET_PATH)/$(BINARY_WINDOWS)

all: test build run

build: goformat
	$(GOBUILD) -o $(BINARY_TARGET_PATH) -v

test:  goformat
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_TARGET_PATH)
	rm -f $(BINARY_TARGET_UNIX_PATH)

run:
	$(GOBUILD) -o $(BINARY_TARGET_PATH) -v
	$(BINARY_TARGET_PATH)

# Cross compilation
build-linux: goformat
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_TARGET_UNIX_PATH) -v

build-windows: goformat
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_TARGET_WINDOWS_PATH) -v

docker-build: check-env-GOPATH goformat
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/ulyssessouza/clf-analyzer-server golang:latest go build -o "$(BINARY_TARGET_UNIX_PATH)" -v

# Util
check-env-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

goformat:
	go fmt .
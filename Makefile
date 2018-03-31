# Configuration
#FIXME Consider using environment variables
APP_PORT=8000
GIN_PORT=3000

# Go parameters
GOCMD=go

GOBUILD_ARGS=-o $(BINARY_TARGET_PATH) -v
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BUILD_TARGET_PATH=dist
BINARY_NAME=clf-anayzer-server
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME)_win

GIN_BUILD_ARGS="$(GOBUILD_ARGS)"

BINARY_TARGET_PATH=$(BUILD_TARGET_PATH)/$(BINARY_NAME)
BINARY_TARGET_UNIX_PATH=$(BUILD_TARGET_PATH)/$(BINARY_UNIX)
BINARY_TARGET_WINDOWS_PATH=$(BUILD_TARGET_PATH)/$(BINARY_WINDOWS)

all: test rundev

build: clean goformat swagger
	$(GOBUILD) $(GOBUILD_ARGS)

test: build
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_TARGET_PATH)

run: build
	$(BINARY_TARGET_PATH)

rundev:
	gin -a $(APP_PORT) -p $(GIN_PORT) --bin $(BINARY_TARGET_PATH) --buildArgs $(GIN_BUILD_ARGS) run main.go

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

swagger:
	swag init
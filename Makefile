PROJECT=kemp-client

BUILD_PATH := $(shell pwd)/.gobuild
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

.PHONY=all get-deps build clean

PROJECT_PATH := "$(BUILD_PATH)/src/github.com/giantswarm"

GOPATH := $(BUILD_PATH)

SOURCE=$(shell find . -name '*.go')

BIN := $(PROJECT)

all: .gobuild get-deps $(BIN)

get-deps: .gobuild
	GOPATH=$(GOPATH) go get -d -v github.com/giantswarm/$(PROJECT)

.gobuild:
	mkdir -p $(PROJECT_PATH)
	cd "$(PROJECT_PATH)" && ln -s ../../../.. $(PROJECT)

$(BIN): $(SOURCE) VERSION
	GOPATH=$(GOPATH) go build -a -ldflags "-X main.projectVersion $(VERSION) -X main.projectBuild $(COMMIT)" -o $(BIN)

clean:
	rm -rf $(BUILD_PATH) $(BIN)

fmt:
	gofmt -w -l .

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

APPS_DIR := $(ROOT_DIR)/apps
API_DIR := $(APPS_DIR)/folder-info
VIEW_DIR := $(APPS_DIR)/folder-view

LINTER_DIR := $(ROOT_DIR)/linter
LINTER_BIN := $(LINTER_DIR)/bin/golangci-lint

ENTRYPOINT := cmd/main.go
COMPILATION_OUTPUT := compiled/main

setup: install-linter

lint: lint-api

build: build-api

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_DIR)/bin v1.57.2

lint-api:
	cd $(API_DIR) && $(LINTER_BIN) run ./internal/...

build-api:
	cd $(API_DIR) && go build -o $(COMPILATION_OUTPUT) $(ENTRYPOINT)

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

APPS_DIR := $(ROOT_DIR)/apps

API_DIR := $(APPS_DIR)/folder-info
API_CONFIGS_DIR := $(API_DIR)/configs

VIEW_DIR := $(APPS_DIR)/folder-view
VIEW_TEMPLATES_DIR := $(VIEW_DIR)/web

LINTER_DIR := $(ROOT_DIR)/linter
LINTER_BIN := $(LINTER_DIR)/bin/golangci-lint

ENTRYPOINT := cmd/main.go
COMPILATION_OUTPUT := compiled/main

setup: install-linter

lint: lint-api lint-view

build: build-api build-view

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_DIR)/bin v1.57.2

run-api: build-api
	cd $(API_DIR) && CONFIG_PATH=$(API_CONFIGS_DIR)/main.json ./$(COMPILATION_OUTPUT)

lint-api:
	cd $(API_DIR) && $(LINTER_BIN) run ./internal/...

build-api:
	cd $(API_DIR) && go build -o $(COMPILATION_OUTPUT) $(ENTRYPOINT)

run-view: build-view
	cd $(VIEW_DIR) && ./$(COMPILATION_OUTPUT) -templates-dir=$(VIEW_TEMPLATES_DIR)

lint-view:
	cd $(VIEW_DIR) && $(LINTER_BIN) run ./internal/...

build-view:
	cd $(VIEW_DIR) && go build -o $(COMPILATION_OUTPUT) $(ENTRYPOINT)

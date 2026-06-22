.DEFAULT_GOAL := help

APP_NAME := go-order-inventory
BIN_DIR := bin
GO ?= go
GOLANGCI_LINT ?= golangci-lint
COMPOSE ?= docker compose
PACKAGES := ./...
TEST_FLAGS ?=
LINT_FLAGS ?=

ifeq ($(OS),Windows_NT)
BINARY := $(BIN_DIR)/$(APP_NAME).exe
else
BINARY := $(BIN_DIR)/$(APP_NAME)
endif

.PHONY: help run dev build clean \
	fmt vet lint tidy mod-download mod-verify \
	test test-verbose test-service test-redis test-all test-race coverage coverage-html \
	check compose-config

help:
	@echo Usage: make target
	@echo Development:
	@echo   run             Run the API locally
	@echo   dev             Start MySQL/Redis, then run the API
	@echo   build           Build the API binary into $(BIN_DIR)/
	@echo   clean           Remove generated build and coverage files
	@echo Quality:
	@echo   fmt             Format all Go packages
	@echo   vet             Run go vet
	@echo   lint            Run golangci-lint - installation required
	@echo   tidy            Update go.mod and go.sum
	@echo   mod-download    Download Go modules
	@echo   mod-verify      Verify downloaded Go modules
	@echo   check           Run format, module verification, vet, and tests
	@echo Tests:
	@echo   test            Run all tests
	@echo   test-verbose    Run all tests with verbose output
	@echo   test-service    Run service tests
	@echo   test-redis      Run Redis integration tests
	@echo   test-all        Run all tests, including Redis integration tests
	@echo   test-race       Run all tests with the race detector
	@echo   coverage        Generate coverage.out
	@echo   coverage-html   Generate coverage.html

run:
	$(GO) run ./cmd

dev: infra-up run

build:
ifeq ($(OS),Windows_NT)
	@if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
else
	mkdir -p "$(BIN_DIR)"
endif
	$(GO) build -trimpath -o "$(BINARY)" ./cmd

clean:
ifeq ($(OS),Windows_NT)
	@if exist "$(subst /,\,$(BIN_DIR))" rmdir /S /Q "$(subst /,\,$(BIN_DIR))"
	@if exist coverage.out del /Q coverage.out
	@if exist coverage.html del /Q coverage.html
else
	rm -rf "$(BIN_DIR)" coverage.out coverage.html
endif

fmt:
	$(GO) fmt $(PACKAGES)

vet:
	$(GO) vet $(PACKAGES)

lint:
ifeq ($(OS),Windows_NT)
	@where "$(GOLANGCI_LINT)" >NUL 2>&1 || (echo golangci-lint is not installed & exit /B 1)
else
	@command -v "$(GOLANGCI_LINT)" >/dev/null 2>&1 || { echo "golangci-lint is not installed"; exit 1; }
endif
	$(GOLANGCI_LINT) run $(LINT_FLAGS) $(PACKAGES)

tidy:
	$(GO) mod tidy

mod-download:
	$(GO) mod download

mod-verify:
	$(GO) mod verify

test:
	$(GO) test $(TEST_FLAGS) $(PACKAGES)

test-verbose:
	$(GO) test -v $(TEST_FLAGS) $(PACKAGES)

test-service:
	$(GO) test -v $(TEST_FLAGS) ./internal/service

test-redis: export RUN_REDIS_TEST := 1
test-redis:
	$(GO) test -v $(TEST_FLAGS) ./internal/bizcache

test-all: test test-redis

test-race:
	$(GO) test -race $(TEST_FLAGS) $(PACKAGES)

coverage:
	$(GO) test $(TEST_FLAGS) -covermode=atomic -coverprofile=coverage.out $(PACKAGES)

coverage-html: coverage
	$(GO) tool cover -html=coverage.out -o coverage.html

check: fmt mod-verify vet test



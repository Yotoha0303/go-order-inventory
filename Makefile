APP_NAME := go-order-inventory
IMAGE_NAME := go-order-inventory:dev
GOLANGCI_LINT ?= golangci-lint

.PHONY: help run test

help:
	@echo Usage: make target
	@echo Targets:
	@echo 	help		Show this help message
	@echo   test        Run go test
	@echo 	run			Run the application locally

test:
	go test ./...

lint:
	$(GOLANGCI_LINT) run ./...
run:
	go run ./cmd
APP_NAME := go-order-inventory
IMAGE_NAME := go-order-inventory:dev


help:
	@echo Usage: make target
	@echo Targets:
	@echo 	help		Show this help message
	@echo 	run			Run the application locally

run:
	go run ./cmd
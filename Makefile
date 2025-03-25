SHELL = /bin/bash

# Load and export environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export $(shell sed 's/=.*//' .env)
else
    $(error .env file not found)
endif

# Get the number of CPU cores for parallelism
#get_cpu_cores := $(shell getconf _NPROCESSORS_ONLN)
# Shell function to determine the number of CPU cores based on the OS
get_cpu_cores = \
  if [ "$$(uname)" = "Linux" ]; then \
    nproc; \
  elif [ "$$(uname)" = "Darwin" ]; then \
    sysctl -n hw.ncpu; \
  else \
    echo "Unsupported OS, default to 1"; \
    echo 1; \
  fi

# Assign the result of the get_cpu_cores shell command to a variable
cpu_cores := $(shell $(get_cpu_cores))

# Project-specific variables
PROJECT_NAME := $(shell basename "$(PWD)" | tr '[:upper:]' '[:lower:]')
VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

DOCKER_IMAGE_NAME := $(PROJECT_NAME):$(VERSION)
DOCKER_CONTAINER_NAME := $(PROJECT_NAME)-$(VERSION)
MODULE := $(shell go list -m)
TEST_COVERAGE_THRESHOLD := 80

# Command to calculate test coverage
testCoverageCmd := $(shell go tool cover -func=coverage.out | grep total | awk '{print $$3}')

# Helper variables
GO_BUILD_CMD := CGO_ENABLED=0 go build $(LDFLAGS) -o $(PROJECT_NAME)
GO_TEST_CMD := go test ./... -v -coverprofile=coverage.out -covermode=count -parallel=$(cpu_cores)

## Start all necessary services and API server
.PHONY: start
start: setup run ## Start all necessary services and API server

## Start only dependencies (Docker containers)
.PHONY: setup
setup: docker-compose-up ## Start only dependencies

## Run the API server
.PHONY: run
run: ## Run the API server
	go run $(LDFLAGS) main.go -version=$(VERSION)

## Start docker-compose services
.PHONY: docker-compose-up
docker-compose-up:
	docker-compose up -d

## Build the API server binary
.PHONY: build
build: ## Build the API server binary
	$(GO_BUILD_CMD) $(MODULE)

## Display the current version of the API server
.PHONY: version
version: ## Display the current version of the API server
	@echo $(VERSION)

## Run tests with coverage
.PHONY: test
test: ## Run tests with coverage
	$(GO_TEST_CMD)

## Generate and display the code coverage report
.PHONY: coverage
coverage: test ## Generate and display the code coverage report
	@echo "Total test coverage:"
	@go tool cover -func=coverage.out | grep total
	@go tool cover -html=coverage.out

## Check if test coverage meets the threshold
.PHONY: ci-coverage
ci-coverage: test ## Check if test coverage meets the threshold
	@coverage=$(shell go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//g'); \
	echo "Current unit test coverage: $$coverage"; \
	echo "Test Coverage Threshold: $(TEST_COVERAGE_THRESHOLD)"; \
	if [ -z "$$coverage" ]; then \
		echo "Test coverage output is empty. Make sure the tests ran successfully."; \
		exit 1; \
	elif [ $$(echo "$$coverage < $(TEST_COVERAGE_THRESHOLD)" | bc) -eq 1 ]; then \
		echo "Test coverage below threshold. Please add more tests."; \
		exit 1; \
	else \
		echo "Test coverage meets the threshold."; \
	fi

## Tidy Go modules
.PHONY: tidy
tidy: ## Tidy Go modules
	go mod tidy

## Format Go code
.PHONY: format
format: ## Format Go code
	go fmt ./...

## Run the linter
.PHONY: lint
lint: ## Run the linter
	golangci-lint run

## Run the linter and fix issues
.PHONY: lint-fix
lint-fix: ## Run the linter and fix issues
	golangci-lint run --fix

## Generate OWASP report
.PHONY: owasp-report
owasp-report: ## Generate OWASP report
	vacuum html-report -z OpenApi-v1.yaml

## Generate Go work file
.PHONY: go-work
go-work: ## Generate Go work file
	go work init .

## Clean all Docker resources
.PHONY: clean
clean: docker-clean ## Clean all Docker resources

## Build the Docker image
.PHONY: docker-build
docker-build: ## Build the Docker image
	$(info ---> Building Docker Image: $(DOCKER_IMAGE_NAME))
	@if [ "$$(uname)" = "Darwin" ]; then \
		DOCKER_BUILDKIT=1 docker-buildx build --output=type=docker --tag $(DOCKER_IMAGE_NAME) --build-arg port=$(port) .; \
	else \
		DOCKER_BUILDKIT=1 docker buildx build --output=type=docker --tag $(DOCKER_IMAGE_NAME) --build-arg port=$(port) .; \
	fi

## Build the Docker image without cache
.PHONY: docker-build-debug
docker-build-debug: ## Build the Docker image without cache
	 $(info ---> Building Docker Image: $(DOCKER_IMAGE_NAME))
	 @if [ "$$(uname)" = "Darwin" ]; then \
	  DOCKER_BUILDKIT=1 docker-buildx build --no-cache --progress=plain --build-arg port=$(port) --tag $(DOCKER_IMAGE_NAME) --output=type=docker .; \
	 else \
	  DOCKER_BUILDKIT=1 docker buildx build --no-cache --progress=plain --build-arg port=$(port) --tag $(DOCKER_IMAGE_NAME) --output=type=docker .; \
	 fi


## Run the Docker container
.PHONY: docker-run
docker-run: ## Run the Docker container
	$(info ---> Running Docker Container: $(DOCKER_CONTAINER_NAME) in Environment: $(profile))
	docker run --name $(DOCKER_CONTAINER_NAME) -it --env environment=$(profile) $(DOCKER_IMAGE_NAME)

## Build and run the Docker container
.PHONY: docker-start
docker-start: docker-build docker-run ## Build and run the Docker container

## Stop the Docker container
.PHONY: docker-stop
docker-stop:
	@if [ -n "$$(docker ps -q --filter name=$(DOCKER_CONTAINER_NAME))" ]; then \
		echo "Stopping container: $(DOCKER_CONTAINER_NAME)"; \
		docker stop $(DOCKER_CONTAINER_NAME); \
	else \
		echo "No container to stop with name: $(DOCKER_CONTAINER_NAME)"; \
	fi

## Remove Docker images and containers
.PHONY: docker-remove
docker-remove:
	@if [ -n "$$(docker ps -a -q --filter name=$(DOCKER_CONTAINER_NAME))" ]; then \
		echo "Removing container: $(DOCKER_CONTAINER_NAME)"; \
		docker rm $(DOCKER_CONTAINER_NAME); \
	else \
		echo "No container to remove with name: $(DOCKER_CONTAINER_NAME)"; \
	fi

	@if [ -n "$$(docker images -q $(DOCKER_IMAGE_NAME))" ]; then \
		echo "Removing image: $(DOCKER_IMAGE_NAME)"; \
		docker rmi $(DOCKER_IMAGE_NAME); \
	else \
		echo "No image to remove with name: $(DOCKER_IMAGE_NAME)"; \
	fi

## Clean all Docker resources
.PHONY: docker-clean
docker-clean: docker-stop docker-remove docker-clean-build-images ## Clean all Docker resources

## Remove build images
.PHONY: docker-clean-build-images
docker-clean-build-images:
	@if [ -n "$$(docker images --filter label="builder=true" -q)" ]; then \
		echo "Removing build images..."; \
		docker rmi $$(docker images --filter label="builder=true" -q); \
	else \
		echo "No build images to remove."; \
	fi

## Display help
.PHONY: help
help:
	@echo
	@echo "Available commands for $(PROJECT_NAME):"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-30s %s\n", $$1, $$2}' | sort
	@echo
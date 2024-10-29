# Load and export environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export $(shell sed 's/=.*//' .env)
else
    $(error .env file not found)
endif

# Project-specific variables
PROJECT_NAME := $(shell basename "$(PWD)" | tr '[:upper:]' '[:lower:]')
VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

DOCKER_IMAGE_NAME := "$(PROJECT_NAME):$(VERSION)"
DOCKER_CONTAINER_NAME := "$(PROJECT_NAME)-$(VERSION)"
MODULE := $(shell go list -m)
TEST_COVERAGE_THRESHOLD = 50

# Test coverage command
testCoverageCmd := $(shell go tool cover -func=coverage.out | grep total | awk '{print $$3}')

# Start all necessary services and API server
.PHONY: start
start: docker-compose-up run

# Start only dependencies
.PHONY: setup
setup: docker-compose-up

# Run the API server
.PHONY: run
run:
	go run $(LDFLAGS) main.go -version=$(VERSION)

# Build the API server binary
.PHONY: build
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(PROJECT_NAME) $(MODULE)

# Display the current version of the API server
.PHONY: version
version:
	@echo $(VERSION)

# Run tests with coverage
.PHONY: test
test:
	go test ./... -v -coverprofile=coverage.out -covermode=count

# Generate and display the code coverage report
.PHONY: coverage
coverage: test
	@go tool cover -func=coverage.out | grep total | awk '{print "Total test coverage: "$$3}'
	@go tool cover -html=coverage.out

# Check if test coverage meets the threshold
.PHONY: ci-coverage
ci-coverage: test
	@echo "Current unit test coverage: $(testCoverageCmd)"
	@echo "Test Coverage Threshold: $(TEST_COVERAGE_THRESHOLD)"
	@if [ $$(echo "$(testCoverageCmd) < $(TEST_COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Failed, please add or update tests to improve test coverage."; \
		exit 1; \
	else \
		echo "OK"; \
	fi

# Tidy Go modules
.PHONY: tidy
tidy:
	go mod tidy

# Format Go code
.PHONY: format
format:
	go fmt ./...

# Run the linter
.PHONY: lint
lint:
	golangci-lint run

# Run the linter and fix issues
.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

# Clean all Docker resources
.PHONY: clean
clean: docker-clean

# Build the Docker image
.PHONY: docker-build
docker-build:
	$(info ---> Building Docker Image: $(DOCKER_IMAGE_NAME))
	docker build -t $(DOCKER_IMAGE_NAME) --build-arg port=$(port) .

# Build the Docker image without cache
.PHONY: docker-build-debug
docker-build-debug:
	$(info ---> Building Docker Image: $(DOCKER_IMAGE_NAME))
	docker build --no-cache --progress=plain -t $(DOCKER_IMAGE_NAME) --build-arg port=$(port) .

# Generate OWASP report
.PHONY: owasp-report
owasp-report:
	vacuum html-report -z OpenApi-v1.yaml

# Generate Go work file
.PHONY: go-work
go-work:
	go work init .

# Run the Docker container
.PHONY: docker-run
docker-run:
	$(info ---> Running Docker Container: $(DOCKER_CONTAINER_NAME) in Environment: $(profile))
	docker run --name $(DOCKER_CONTAINER_NAME) -it --env environment=$(profile) $(DOCKER_IMAGE_NAME)

# Build and run the Docker container
.PHONY: docker-start
docker-start: docker-build docker-run

# Stop the Docker container
.PHONY: docker-stop
docker-stop:
	docker stop $(DOCKER_CONTAINER_NAME)

# Remove Docker images and containers
.PHONY: docker-remove
docker-remove:
	docker rm $(DOCKER_CONTAINER_NAME)
	docker rmi $(DOCKER_IMAGE_NAME)

# Clean all Docker resources
.PHONY: docker-clean
docker-clean: docker-clean-service-images docker-clean-build-images

# Stop and remove service images
.PHONY: docker-clean-service-images
docker-clean-service-images: docker-stop docker-remove

# Remove build images
.PHONY: docker-clean-build-images
docker-clean-build-images:
	docker rmi $$(docker images --filter label="builder=true" -q)

# Display help
.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command to run in $(PROJECT_NAME):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo

-include .env

export environment=local

PROJECT_NAME := $(shell basename "$(PWD)" | tr '[:upper:]' '[:lower:]')

# GIT commit id will be used as version of the application
VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags "-X main.version=${VERSION}"

DOCKER_IMAGE_NAME := "$(PROJECT_NAME):$(VERSION)"
DOCKER_CONTAINER_NAME := "$(PROJECT_NAME)-$(VERSION)"

MODULE = $(shell go list -m)

testCoverageCmd := $(shell go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+')
TEST_COVERAGE_THRESHOLD = 50

## start: Starts everything that is required to serve the APIs
start:
	docker-compose up -d
	make run

## run: Run the API server alone (without supplementary services such as DB etc.,)
run:
	go run ${LDFLAGS} main.go -version="${VERSION}"

## build: Build the API server binary
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o ${PROJECT_NAME} $(MODULE)

## version: Display the current version of the API server
version:
	@echo $(VERSION)

## test: Run tests
test:
	go test ./... -v -coverprofile coverage.out -covermode count

## coverage: Measures and generate code coverage report
coverage:
	@go test ./... -v -coverprofile coverage.out -covermode count
	@go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+' | xargs -I {} echo "Total test coverage: {}%"
	@go tool cover -html=coverage.out

ci-coverage:
	@echo "Current unit test coverage: $(testCoverageCmd)"
	@echo "Test Coverage Threshold: $(TEST_COVERAGE_THRESHOLD)"
	@echo "-----------------------"

	@if [ "$(shell echo "$(testCoverageCmd) < $(TEST_COVERAGE_THRESHOLD)" | bc -l)" -eq 1 ]; then \
		echo "Failed, Please add or update tests to improve test code coverage."; \
		exit 1; \
	else \
		echo "OK"; \
	fi

## tidy: Tidy go modules
tidy:
	go mod tidy

## format: Format go code
format:
	go fmt ./...

## lint: Run linter
lint:
	golangci-lint run

## lint-fix: Run linter and fix the issues
lint-fix:
	golangci-lint run --fix

## docker-build: Build the API server as a docker image
docker-build:
	$(info ---> Building Docker Image: ${DOCKER_IMAGE_NAME}, Exposed Port: ${port})
	docker build -t ${DOCKER_IMAGE_NAME} . \
		--build-arg port=${port} \

docker-build-debug:
	$(info ---> Building Docker Image: ${DOCKER_IMAGE_NAME}, Exposed Port: ${port})
	docker build --no-cache --progress plain -t ${DOCKER_IMAGE_NAME} . \
		--build-arg port=${port} \

## docker-run: Run the API server as a docker container
docker-run:
	$(info ---> Running Docker Container: ${DOCKER_CONTAINER_NAME} in Environment: ${profile})
	docker run --name ${DOCKER_CONTAINER_NAME} -it \
				--env environment=${profile} \
				$(DOCKER_IMAGE_NAME)

## docker-start: Builds Docker image and runs it.
docker-start: docker-build docker-run

## docker-stop: Stops the docker container
docker-stop:
	docker stop $(DOCKER_CONTAINER_NAME)

## docker-remove: Removes the docker images and containers	
docker-remove:
	docker rm $(DOCKER_CONTAINER_NAME)
	docker rmi $(DOCKER_IMAGE_NAME)

## docker-clean: Cleans all docker resources
docker-clean: docker-clean-service-images docker-clean-build-images

## docker-clean-service-images: Stops and Removes the service images
docker-clean-service-images: docker-stop docker-remove

## docker-clean-build-images: Removes build images
docker-clean-build-images: 
	docker rmi $(docker images --filter label="builder=true")

## owasp-report: Generate OWASP report
owasp-report:
	vacuum html-report -z OpenApi-v1.yaml

## go-work: Generate the go work file
go-work:
	go work init .

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |	sed -e 's/^/ /'
	@echo
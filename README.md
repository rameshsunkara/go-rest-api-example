[![Build Status](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml/badge.svg)](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml?query=+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/rameshsunkara/go-rest-api-example)](https://goreportcard.com/report/github.com/rameshsunkara/go-rest-api-example)
[![codecov](https://codecov.io/gh/rameshsunkara/go-rest-api-example/branch/main/graph/badge.svg)](https://app.codecov.io/gh/rameshsunkara/go-rest-api-example)


# REST API microservice in Go

## [Why this ?](#why-this--1)

## Offered Features

### API Features:
1. **OWASP Compliant Open API 3 Specification**: Refer to [OpenApi-v1.yaml](./OpenApi-v1.yaml) for details.
2. **Middleware for**:
   - **Logging**: Facilitates debugging and monitoring processes.
   - **Authentication**: Provides a placeholder for diverse authentication mechanisms.
   - **Tracing by Request ID**: Assists in debugging procedures.
   - **Panic Recovery**: Ensures service continuity by managing unexpected errors.
   - **Common Security Headers**: Safeguards the service against potential vulnerabilities.
   - **Query Parameters Validation**: Enhances service security by validating query parameters.
3. **Standardized Error Handling**: All errors are managed and returned in a uniform format.
4. **Versioning**
5. **Model Management**:
   - Internally used data models differ from those exposed to clients, ensuring separation and security.

### Go Application Features:
1. **Configuration Management**: via Environment Variables
2. **Dockerized Environment**: Facilitates service deployment.
3. **Makefile**: Automates common tasks for developers.
4. **Git Action**: Automates build processes, runs tests, and generates code coverage.
5. **Integrated Go Formatter and Linter**: Promotes code quality and consistency.
6. **Secrets Loading Mechanism from Sidecar**
7. **Support for Multiple Databases**: Enables connections to various database systems.
8. **Best Practices for MongoDB Connection**
9. **Effective Mocking Practices for Unit Test Patterns**
10. **Seed Data**: for Local Development
11. **Standardized Filename Conventions**: Enhances code readability.
12. **Multi-Stage Docker Build**: Accelerates build processes.
13. **Versioning** Utilizing Git Commit History

## Folder Structure

```
go-rest-api-example/
├── main.go
├── internal/
│   ├── db
│   ├── errors
│   ├── handlers
│   ├── logger
│   ├── middleware
│   ├── models
│   ├── server
│   ├── util
│   └── mockData
├── localDevelopment/
├── Makefile
├── Dockerfile
├── OpenApi-vi.yaml
├── docker-compose.yaml
└── OpenApi-v1.postman_collection.json
```

## QuickStart

Pre-requisites: Docker, Docker Compose, Make

1. Start the service

        make start

Other Options:

   Choose a command to run in go-rest-api-example:
   
      start                         Starts everything that is required to serve the APIs
      run                           Run the API server alone (without supplementary services such as DB etc.,)
      build                         Build the API server binary
      version                       Display the current version of the API server
      test                          Run tests
      coverage                      Measures and generate code coverage report
      tidy                          Tidy go modules
      format                        Format go code
      lint                          Run linter
      lint-fix                      Run linter and fix the issues
      docker-build                  Build the API server as a docker image
      docker-run                    Run the API server as a docker container
      docker-start                  Builds Docker image and runs it.
      docker-stop                   Stops the docker container
      docker-remove                 Removes the docker images and containers        
      docker-clean                  Cleans all docker resources
      docker-clean-service-images   Stops and Removes the service images
      docker-clean-build-images     Removes build images
      owasp-report                  Generate OWASP report
      go-work                       Generate the go work file

## Tools

1. Routing - [Gin](https://github.com/gin-gonic/gin)
2. Logging - [zerolog](https://github.com/rs/zerolog)
3. Database - [MongoDB](https://www.mongodb.com/)
4. Container - [Docker](https://www.docker.com/)

## TODO

-  Add more and clear documentation about the features this offers and how to replace tools
-  Add DB Migration Support
-  Add more environmental profiles and obey all [12-Factor App rules](https://12factor.net/ru/)
-  Implement all OWASP security checks specified in the API Spec
-  Improve error codes and messages
-  Add git hooks for pre-commit and pre-push

## Good to have

-  Improve data model and add more fields
-  Deploy to cloud
-  Implement Update Operations mentioned in the API Spec

## References

- [gin-boilerplate](https://github.com/Massad/gin-boilerplate)
- [go-rest-api](https://github.com/qiangxue/go-rest-api)
- [go-base](https://github.com/dhax/go-base)

## Contribute

- Please feel free to Open PRs
- Please create issues with any problem you noticed
- Please suggest any improvements

## Why this ?


I embarked on the endeavor of crafting my own open-source boilerplate repository for several reasons:

After years of developing Full Stack applications utilizing ReactJS and JVM-based languages, I observed that existing boilerplates tended to be either excessive or insufficient for my needs. Consequently, I resolved to construct my own, while adhering rigorously to the principles and guidelines of Go. While similarities with popular Go boilerplate templates may be evident, I have customized this repository to better align with my preferences and accumulated experiences. (My apologies if I inadvertently overlooked crediting any existing templates.)

I yearned for the autonomy to meticulously select the tools for fundamental functionalities such as Routing, Logging, and Configuration Management, ensuring seamless alignment with my personal preferences and specific requirements.

## What this is not ?

- This isn't a complete solution for all your needs. It's more like a basic template to kickstart your project.
- This isn't the best place to begin if you want to make an online store. What I've provided is just a simple tool for managing data through an API.
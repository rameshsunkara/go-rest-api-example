[![Build Status](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml/badge.svg)](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml?query=+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/rameshsunkara/go-rest-api-example)](https://goreportcard.com/report/github.com/rameshsunkara/go-rest-api-example)
[![codecov](https://codecov.io/gh/rameshsunkara/go-rest-api-example/branch/main/graph/badge.svg)](https://app.codecov.io/gh/rameshsunkara/go-rest-api-example)


# REST API microservice in golang

## Why?

There are many open source boilerplate repos but why I did this ?

1. Coming from years of building Full Stack application in ReactJS and JVM based languages, I did not like any of them
   completely.
   So I created my own while obeying 'GO' principles and guidelines.
   You will find a lot of similarities in this repo when compared to the most popular go boilerplate templates because I
   probably borrowed ideas from them. (my apologies if I failed to miss any of them in the references)

2. I wanted to pick the tools for Routing, Logging, Configuration Management etc., to my liking and preferences.

3. I wanted a version where I have full control to change/update based on my professional work requirements.

### QuickStart

Pre-requisites: Docker, Docker Compose, Make

1. Start the service

        make start

   If you are a fan of Postman, import the included [Postman collection](orders.postman_collection.json) or use the [OpenAPI3 Spec file](./OpenApi-v1.yaml).

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



### Tools

1. Routing - [Gin](https://github.com/gin-gonic/gin)
1. Logging - [zerolog](https://github.com/rs/zerolog)
1. Database - [Mongo](https://www.mongodb.com/)
1. Container - [Docker](https://www.docker.com/)

### Features

- OpenApi3.1 Spec
- Easy to use 'make' tasks to do everything
- Multi-Stage container build (cache enabled)
- Versioning using git commit (both Application and Docker objects)
- Git Actions to build, security analysis and to run code coverage
- Templated Docker and Make files

### TODO

- [ ] Add more and clear documentation about the features this offers and how to replace tools
- [ ] Add DB Migration Support
- [ ] Add more profiles and obey all [12-Factor App rules](https://12factor.net/ru/)
- [ ] Deploy to cloud
- [ ] Add missing references/inspirations

### References

- [gin-boilerplate](https://github.com/Massad/gin-boilerplate)
- [go-rest-api](https://github.com/qiangxue/go-rest-api)
- [go-base](https://github.com/dhax/go-base)

### Contribute

- Please feel free to Open PRs
- Please create issues with any problem you noticed
- Please suggest any improvements

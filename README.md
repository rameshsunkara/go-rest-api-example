[![Build Status](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml/badge.svg)](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml?query=+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/rameshsunkara/go-rest-api-example)](https://goreportcard.com/report/github.com/rameshsunkara/go-rest-api-example)
[![codecov](https://codecov.io/gh/rameshsunkara/go-rest-api-example/branch/main/graph/badge.svg)](https://app.codecov.io/gh/rameshsunkara/go-rest-api-example)


# REST API microservice in Go

## [Why this ?](#why-this--1)

## What does it offer ?

### API Features:
1. OWASP Compliant [Open API 3 Spec](./OpenApi-v1.yaml)
2. Middleware for 
   - Logging : Helps in debugging and monitoring
   - Authentication : Placeholder for different authentication mechanisms
   - Tracing by Request ID : Helps in debugging
   - Panic Recovery : Helps in keeping the service up
   - Common Security Headers : Keeps the service secure
   - Query Params Validation : Helps in keeping the service secure
3. Standard Error Handling
   - All errors are handled and returned in a standard format
4. Versioning
5. Model Management
   - Generally, the data model used internally is different from the data model exposed to the client.
     This helps in keeping the internal model separate from the exposed model.
 
### Go Application Features:
1. Configuration Management through Environment Variables
2. A Dockerized environment to run the service
3. A Makefile to do all common tasks
4. A Git Action to build, run tests, generate code coverage
5. Integrated GO Formatter and Linter
6. Mechanism to load secrets from Sidecar
7. Enables connecting to multiple databases
8. Follows the best practices for connecting to MongoDB
9. Good mocking practises for Unit test patterns
10. Seed data for local development
11. Standard filename conventions for better readability
12. Multi-Stage Docker build for faster builds
13. Versioning using git commit


### QuickStart

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

### Tools

1. Routing - [Gin](https://github.com/gin-gonic/gin)
2. Logging - [zerolog](https://github.com/rs/zerolog)
3. Database - [MongoDB](https://www.mongodb.com/)
4. Container - [Docker](https://www.docker.com/)

### TODO

- [ ] Add more and clear documentation about the features this offers and how to replace tools
- [ ] Add DB Migration Support
- [ ] Add more profiles and obey all [12-Factor App rules](https://12factor.net/ru/)
- [ ] Deploy to cloud
- [ ] Implement all OWASP security checks specified in the API Spec
- [ ] Improve error codes and messages
- [ ] Add missing references/inspirations
- [ ] Implement Update Operations mentioned in the API Spec
- [ ] Improve data model and add more fields
- [ ] Add git hooks for pre-commit and pre-push

### References

- [gin-boilerplate](https://github.com/Massad/gin-boilerplate)
- [go-rest-api](https://github.com/qiangxue/go-rest-api)
- [go-base](https://github.com/dhax/go-base)

### Contribute

- Please feel free to Open PRs
- Please create issues with any problem you noticed
- Please suggest any improvements

#### Why this ?

I ventured into creating my own open-source boilerplate repository for several reasons:

1. After years of crafting Full Stack applications using ReactJS and JVM-based languages, I found existing boilerplate's are either too much or too little. 
   Hence, I decided to develop my own while adhering to the principles and guidelines of Go. 
   While you may notice similarities with popular Go boilerplate templates, I've tailored this repository to align more closely with my preferences and experiences. 
   (Apologies if I inadvertently missed crediting any existing templates.)

2. I desired the freedom to handpick the tools for essential functionalities like Routing, Logging, and Configuration Management, ensuring they align perfectly with my preferences and requirements.

3. Creating my own version allows me complete control to adapt and update the boilerplate according to the specific needs and demands of my professional work. This flexibility enables me to continually refine and optimize the repository based on evolving project requirements.
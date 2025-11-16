# REST API microservice in Go

[![Build Status](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml/badge.svg)](https://github.com/rameshsunkara/go-rest-api-example/actions/workflows/cibuild.yml?query=+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/rameshsunkara/go-rest-api-example)](https://goreportcard.com/report/github.com/rameshsunkara/go-rest-api-example)
[![codecov](https://codecov.io/gh/rameshsunkara/go-rest-api-example/branch/main/graph/badge.svg)](https://app.codecov.io/gh/rameshsunkara/go-rest-api-example)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)

> A production-ready REST API boilerplate built with Go, featuring MongoDB integration, comprehensive middleware, flight recorder tracing, and modern development practices.

![Go REST Api](go-rest-api.svg)

## ✨ Highlights

- 🚀 **Production-Ready**: Graceful shutdown, health checks, structured logging
- 🔒 **Security-First**: OWASP compliant, multi-tier auth, security headers
- 📊 **Observability**: Flight recorder tracing, Prometheus metrics, pprof profiling
- 🧪 **Test Coverage**: 70%+ coverage threshold with parallel testing
- 🐳 **Docker-Ready**: Multi-stage builds with BuildKit optimization
- 📝 **Well-Documented**: OpenAPI 3 specification with Postman collection

## 🚀 Quick Start

```bash
# Clone and start
git clone https://github.com/rameshsunkara/go-rest-api-example.git
cd go-rest-api-example
make start

# Your API is now running at http://localhost:8080
curl http://localhost:8080/healthz
```

## 📋 Table of Contents

- [Features](#-key-features)
- [Architecture](#️-architecture)
- [Getting Started](#-getting-started)
- [Available Commands](#-available-commands)
- [Tools & Stack](#-tools--stack)
- [Contributing](#contribute)

## 🎯 Key Features

### API Features

1. **OWASP Compliant Open API 3 Specification**: Refer to [OpenApi-v1.yaml](./OpenApi-v1.yaml) for details.
2. **Production-Ready Health Checks**:
   - `/healthz` endpoint with proper HTTP status codes (204/424)
   - Database connectivity validation
   - Dependency health monitoring
3. **Comprehensive Middleware Stack**:
   - **Request Logging**: Structured logging with request correlation
   - **Authentication**: Multi-tier auth (external/internal APIs)
   - **Request ID Tracing**: End-to-end request tracking
   - **Panic Recovery**: Graceful error handling and recovery
   - **Security Headers**: OWASP-compliant security header injection
   - **Query Validation**: Input validation and sanitization
   - **Compression**: Automatic response compression (gzip)
4. **Flight Recorder Integration**: Automatic trace capture for slow requests using Go 1.25's built-in flight recorder.
5. **Standardized Error Handling**: Consistent error response format across all endpoints
6. **API Versioning**: URL-based versioning with backward compatibility
7. **Internal vs External APIs**: Separate authentication and access controls
8. **Model Separation**: Clear distinction between internal and external data representations

### Go Application Features

1. **Configuration Management**: Environment-based configuration with validation
2. **Graceful Shutdown**: Proper signal handling with resource cleanup and connection draining
3. **Production-Ready MongoDB Integration**:
   - Connection pooling and health checks
   - Functional options pattern for flexible configuration
   - SRV and replica set support
   - Credential management via sidecar files
   - Query logging for debugging
4. **Comprehensive Health Checks**: `/healthz` endpoint with database connectivity validation
5. **Structured Logging**: Zero-allocation JSON logging with request tracing
6. **Secrets Management**: Secure credential loading from sidecar files
7. **Effective Mocking**: Interface-based design enabling comprehensive unit testing
8. **Database Indexing**: Automatic index creation for optimal query performance
9. **Idiomatic Go Architecture**: Clean separation of concerns with dependency injection
10. **Parallel Testing**: Race condition detection with atomic coverage reporting
11. **Context-Aware Operations**: Proper context propagation for cancellation and timeouts
12. **Resource Management**: Automatic cleanup of connections and resources

### Tooling

1. **Dockerized Environment**: Facilitates service deployment using DOCKER_BUILDKIT.
2. **Makefile**: Automates common tasks for developers.
3. **GitHub Actions**: Automates building, testing, code coverage reporting, and enforces the required test coverage threshold.
4. **Multi-Stage Docker Build**: Accelerates build processes.

## 🏗️ Architecture

### 📁 Folder Structure

```text
go-rest-api-example/
├── main.go
├── internal/           # Private application code
│   ├── config/         # Configuration management
│   ├── db/             # Database repositories and data access
│   ├── errors/         # Application error definitions
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware components
│   ├── models/         # Domain models and data structures
│   ├── server/         # HTTP server setup and lifecycle
│   ├── utilities/      # Internal utilities
│   └── mockData/       # Test and development data
├── pkg/                # Public packages (can be imported)
│   ├── logger/         # Structured logging utilities
│   └── mongodb/        # MongoDB connection management
├── localDevelopment/   # Local dev setup (DB init scripts, etc.)
├── Makefile            # Development automation
├── Dockerfile          # Container image definition
├── docker-compose.yaml # Local development services
├── OpenApi-v1.yaml     # API specification
└── OpenApi-v1.postman_collection.json
```

### ➡️ Control Flow

```mermaid
flowchart LR
   Request e1@==> Server
   e1@{ animate: true }
   Server e2@==> Router
   e2@{ animate: true }
   M@{ shape: processes, label: "Middlewares" }
   Router e3@==> M
   e3@{ animate: true }
   C@{ shape: processes, label: "Handlers" }
   M e4@==> C
   e4@{ animate: true }
   R@{ shape: processes, label: "Repos(DAO)" }
   C e5@==> R
   e5@{ animate: true }
   id1[(Database)]
   R e6@==> id1
   e6@{ animate: true }
```

1. **Request**: Server receives the incoming request.
2. **Server**: Server processes the request and forwards it to the router.
3. **Router**: Router directs the request to the appropriate middleware(s).
4. **Middlewares**: The middlewares handle various tasks such as logging, authentication, security headers, tracing etc.,
5. **Handlers**: The request is passed to the appropriate handler, which validates the request and forwards it to the repository layer.
6. **Repos(DAO)**: The repository layer communicates with the database to perform CRUD operations.

## 🚀 Getting Started

### Prerequisites

- Docker and [Docker Compose](https://docs.docker.com/compose/install/)
- Make

### Start the Application

```bash
git clone https://github.com/rameshsunkara/go-rest-api-example.git
cd go-rest-api-example
make start
```

Your API is now running at `http://localhost:8080`

**Try it out:**
```bash
curl http://localhost:8080/api/v1/healthz
curl http://localhost:8080/api/v1/orders
```

## 📟 Available Commands

### Essential Commands

```makefile
start                          Start all necessary services and API server
run                            Run the API server (requires dependencies running)
setup                          Start only dependencies (MongoDB)
test                           Run tests with coverage
```

### Development Commands

```makefile
lint                           Run the linter
lint-fix                       Run the linter and fix issues
trace                          Analyze a trace file (usage: make trace TRACE_FILE=./traces/slow-request-GET-orders-1234567890.trace)
clean                          Clean all Docker resources (keeps database data)
clean-all                      Clean all Docker resources including volumes (removes database data)
coverage                       Generate and display the code coverage report
```

### CI/CD Commands

```makefile
build                          Build the API server binary
ci-coverage                    Check if test coverage meets the threshold
format                         Format Go code
version                        Display the current version of the API server
```

### Docker Commands

```makefile
docker-build                   Build the Docker image
docker-start                   Build and run the Docker container
docker-clean                   Clean all Docker resources
```

> 💡 **Tip**: Run `make help` to see all available commands.

### Additional Prerequisites for Development

- [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation) - For linting
- [docker-buildx](https://docs.docker.com/buildx/working-with-buildx/) - For multi-platform builds

## 🛠 Tools & Stack

| Category | Technology |
|----------|-----------|
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Logging** | [zerolog](https://github.com/rs/zerolog) |
| **Database** | [MongoDB](https://www.mongodb.com/) |
| **Container** | [Docker](https://www.docker.com/) + BuildKit |
| **Tracing** | Go 1.25 Flight Recorder |
| **Profiling** | [pprof](https://golang.org/pkg/net/http/pprof/) |

## 📚 Additional Resources

### Roadmap

<details>
<summary>Click to expand planned features</summary>

- [ ] Add comprehensive API documentation with examples
- [ ] Implement database migration system
- [ ] Add distributed tracing (OpenTelemetry integration)
- [ ] Add metrics collection and Prometheus integration
- [ ] Add git hooks for pre-commit and pre-push
- [ ] Implement all remaining OWASP security checks

</details>

### Nice to Have

<details>
<summary>Future enhancements</summary>

- **Enhanced Data Models**: Add validation, relationships, and business logic
- **Cloud Deployment**: Kubernetes manifests and Helm charts
- **Advanced Monitoring**: APM integration, alerting, and dashboards
- **Caching Layer**: Redis integration for performance optimization
- **Multi-database Support**: PostgreSQL, CockroachDB adapters
- **Performance Testing**: Load testing scenarios and benchmarks

</details>

### References

- [gin-boilerplate](https://github.com/Massad/gin-boilerplate)
- [go-rest-api](https://github.com/qiangxue/go-rest-api)
- [go-base](https://github.com/dhax/go-base)

## 🤝 Contribute

Contributions are welcome! Here's how you can help:

- 🐛 **Found a bug?** [Open an issue](https://github.com/rameshsunkara/go-rest-api-example/issues)
- 💡 **Have a feature idea?** [Start a discussion](https://github.com/rameshsunkara/go-rest-api-example/discussions)
- 🔧 **Want to contribute code?** Fork the repo and submit a PR

## 📖 Why This Project?

After years of developing Full Stack applications using ReactJS and JVM-based languages, I found existing Go boilerplates were either too opinionated or too minimal. 

In the age of Copilot and Cursor, you might wonder: "Why another REST API boilerplate?"
Fair question. But this isn't just another "Hello World" REST API—this is about building for **enterprise**.

### What Makes It Different

This project strikes a balance between simplicity and production-readiness:

✅ **Just Right**: Not too bloated, not too minimal  
✅ **Best Practices**: Follows Go idioms and patterns  
✅ **Production-Tested**: Battle-tested patterns from real-world applications  
✅ **Flexible**: Easy to customize for your specific needs  

### Beyond Performance: Building for Enterprise

In enterprise software, it's not always about raw performance alone. It's about building systems that:

- **Scale with your team**: Code that's readable by engineers at all skill levels, not just the most experienced
- **Enable safe changes**: Clear boundaries between components so changes don't cascade unexpectedly, with each layer independently testable
- **Prevent bad practices**: Built-in guardrails that guide developers toward correct patterns
- **Simplify troubleshooting**: When issues hit production (and they will), you can diagnose them quickly
- **Measure what matters**: Easy performance profiling and metrics collection without major refactoring
- **Support evolution**: Adding new features doesn't require rewriting existing code
- **Delight developers**: New team members become productive in minutes, not days

### What This Is NOT

❌ A complete e-commerce solution
❌ A framework that does everything for you
❌ The only way to structure a Go API

**This is a solid foundation to build upon.** Take what you need, leave what you don't.

---

<div align="center">

**⭐ If you find this helpful, please consider giving it a star! ⭐**

</div>
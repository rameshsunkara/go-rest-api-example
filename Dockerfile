# syntax=docker/dockerfile:1.4  # Enable BuildKit features

# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder
LABEL stage=builder

# Set working directory
WORKDIR /app

# Install dependencies (using cache for Go modules)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download && go mod verify

# Copy source files and build the binary (using cache for build artifacts)
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    make build

# Stage 2: Create a minimal final image
FROM gcr.io/distroless/static-debian11
COPY --from=builder /go-rest-api-example /app/go-rest-api-example
ENTRYPOINT ["/app/go-rest-api-example"]
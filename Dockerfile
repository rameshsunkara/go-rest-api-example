# syntax=docker/dockerfile:1.4  # Enable BuildKit features

# Stage 1: Build the Go binary
FROM golang:1.24 AS builder
LABEL stage=builder

# Set working directory
WORKDIR /app

# Install dependencies (using cache for Go modules)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download && go mod verify

# Copy source files and build the binary (using cache for build artifacts)
# Since the working directory is set as app, the executable generated will be named as app too
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    make build

# Stage 2: Create a minimal final image
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]
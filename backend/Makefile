.PHONY: build run-api run-scheduler run-all clean test test-verbose test-coverage

# Build the application
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/scheduler cmd/worker/scheduler/main.go

# Run the API server
run-api:
	go run cmd/api/main.go

# Run the scheduler worker
run-scheduler:
	go run cmd/worker/scheduler/main.go

# Run both API and scheduler (in separate terminals)
run-all:
	@echo "Starting API server..."
	@go run cmd/api/main.go & \
	echo "Starting scheduler..." && \
	go run cmd/worker/scheduler/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests (quiet)
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with short mode
test-short:
	go test -v -short ./...

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/scheduler ./cmd/worker/scheduler

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run-api       - Run the API server"
	@echo "  run-scheduler - Run the scheduler worker"
	@echo "  run-all       - Run both API and scheduler"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests (quiet)"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-short    - Run tests in short mode"
	@echo "  build-prod    - Build for production"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
.PHONY: build run-api run-scheduler run-all clean test test-verbose test-coverage setup setup-env setup-auth setup-db migrate

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/scheduler cmd/worker/scheduler/main.go

run-api:
	go run cmd/api/main.go

run-scheduler:
	go run cmd/worker/scheduler/main.go

run-all:
	@echo "Starting API server..."
	@go run cmd/api/main.go & \
	echo "Starting scheduler..." && \
	go run cmd/worker/scheduler/main.go

clean:
	rm -rf bin/

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-short:
	go test -v -short ./...

build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/scheduler ./cmd/worker/scheduler

deps:
	go mod tidy
	go mod download

fmt:
	go fmt ./...

lint:
	golangci-lint run

setup-env:
	@echo "Setting up .env file..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file from template..."; \
		echo "# Database Configuration" > .env; \
		echo "DATABASE_URL=\"postgresql://username:password@localhost:26257/stock_insights?sslmode=disable\"" >> .env; \
		echo "" >> .env; \
		echo "# External API Configuration" >> .env; \
		echo "EXTERNAL_API_URL=\"https://api.karenai.click\"" >> .env; \
		echo "EXTERNAL_API_KEY=\"your_external_api_key_here\"" >> .env; \
		echo "" >> .env; \
		echo "# Application Configuration" >> .env; \
		echo "PORT=\"8080\"" >> .env; \
		echo "ENVIRONMENT=\"development\"" >> .env; \
		echo "LOG_LEVEL=\"info\"" >> .env; \
		echo "RATE_LIMIT=\"100\"" >> .env; \
		echo "CACHE_TTL=\"5m\"" >> .env; \
		echo "" >> .env; \
		echo "# Admin Authentication (will be generated)" >> .env; \
		echo "JWT_SECRET=\"\"" >> .env; \
		echo "ADMIN_TOKEN=\"\"" >> .env; \
		echo "" >> .env; \
		echo "✅ .env file created!"; \
		echo "📝 Please edit .env file with your database URL and external API key"; \
	else \
		echo "✅ .env file already exists"; \
	fi

setup-auth:
	@echo "Setting up authentication..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Run 'make setup-env' first"; \
		exit 1; \
	fi
	@go run cmd/setup-auth/main.go
	@echo "✅ Authentication setup completed!"

setup-db:
	@echo "Setting up database..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Run 'make setup-env' first"; \
		exit 1; \
	fi
	@go run cmd/migrate/main.go
	@echo "✅ Database setup completed!"

migrate:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go

setup: setup-env setup-auth setup-db deps
	@echo ""
	@echo "🎉 Setup completed successfully!"
	@echo ""
	@echo "📋 Next steps:"
	@echo "1. Edit .env file with your database URL and external API key"
	@echo "2. Start your database (CockroachDB or PostgreSQL)"
	@echo "3. Run 'make run-api' to start the server"
	@echo ""
	@echo "🔗 Quick start:"
	@echo "   make run-api"
	@echo ""
	@echo "🧪 Test the setup:"
	@echo "   curl http://localhost:8080/api/v1/public/health"

start: deps run-api

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
	@echo ""
	@echo "Setup commands:"
	@echo "  setup         - Complete setup (env + auth + db + deps)"
	@echo "  setup-env     - Create .env file"
	@echo "  setup-auth    - Generate JWT secret and admin token"
	@echo "  setup-db      - Run database migrations"
	@echo "  migrate       - Run database migrations"
	@echo "  start         - Quick start (deps + run-api)"
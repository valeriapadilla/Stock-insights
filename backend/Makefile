.PHONY: test test-coverage test-short build

test:
	go test ./... -v

test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-short:
	go test ./... -v -short

build:
	go build -o main cmd/api/main.go
.PHONY: run build test test-e2e clean tidy

# Run the application
run:
	go run ./cmd/api

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run all tests
test:
	go test ./... -v

# Run e2e tests only
test-e2e:
	go test ./test/e2e/... -v

# Clean build artifacts and test files
clean:
	rm -rf bin/
	rm -f test.db

# Tidy dependencies
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Help
help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run all tests"
	@echo "  make test-e2e  - Run e2e tests only"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make tidy      - Tidy go modules"
	@echo "  make deps      - Download dependencies"

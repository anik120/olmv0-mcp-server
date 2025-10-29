BINARY_NAME=olmv0-mcp-server
BUILD_DIR=bin
MAIN_PATH=cmd/olmv0-mcp-server

.PHONY: all build clean test run help

all: build

## build: Build the olmv0-mcp-server binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## run: Build and run the server
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./$(MAIN_PATH)

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

## lint: Run linting
lint:
	@echo "Running linter..."
	@golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## mod: Update dependencies
mod:
	@echo "Updating dependencies..."
	@go mod tidy
	@go mod vendor

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t olmv0-mcp-server:latest .

## docker-run: Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 olmv0-mcp-server:latest

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sort
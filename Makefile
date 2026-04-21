# Makefile для SOCKS5 прокси

BINARY_NAME := socks5-proxy
BUILD_DIR := build
MODULE_NAME := go-socks5-relay

.PHONY: init
init:
	@echo "Initializing Go module..."
	@if [ ! -f go.mod ]; then \
		go mod init $(MODULE_NAME); \
		echo "Module initialized: $(MODULE_NAME)"; \
	else \
		echo "go.mod already exists"; \
	fi

.PHONY: deps
deps: init
	@echo "Downloading dependencies..."
	go mod tidy
	go mod verify
	@echo "Dependencies ready"

.PHONY: generate-env
generate-env:
	@echo "Generating config/.env file with random credentials..."
	@./scripts/generate-env.sh config/.env.example config/.env

.PHONY: generate-compose-env
generate-compose-env:
	@echo "Generating .env for docker-compose..."
	@./scripts/generate-compose-env.sh

.PHONY: setup
setup: generate-env
	@echo "Setup complete. .env file ready."

.PHONY: build
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/socks5-proxy
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: run
run: generate-env build
	@echo "Starting proxy..."
	$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run-debug
run-debug: generate-env build
	@echo "Starting proxy in debug mode..."
	$(BUILD_DIR)/$(BINARY_NAME) -debug

.PHONY: run-dev
run-dev: generate-env build
	@echo "Starting proxy in development mode (debug + verbose logs)..."
	$(BUILD_DIR)/$(BINARY_NAME) -log-level debug

.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f go.mod go.sum
	@echo "Clean complete"

.PHONY: test
test: deps
	go test -v ./...

.PHONY: docker-up
docker-up:
	@[ -f config/.env ] || ./scripts/generate-env.sh config/.env.example config/.env
	@echo "Starting docker-compose..."
	docker compose --env-file config/.env up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping docker-compose..."
	docker-compose down

.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	docker-compose build

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Build and run"
	@echo "  make run-debug    - Build and run with debug mode"
	@echo "  make run-dev      - Build and run with development settings"
	@echo "  make clean        - Remove build artifacts and module files"
	@echo "  make deps         - Download dependencies"
	@echo "  make test         - Run tests"
	@echo "  make setup        - Generate .env file with random credentials"
	@echo "  make generate-env - Generate .env file only"
	@echo "  make generate-compose-env - Generate .env for docker-compose"
	@echo "  make docker-up    - Start docker-compose"
	@echo "  make docker-down  - Stop docker-compose"
	@echo "  make docker-build - Build docker image"
	@echo ""
	@echo "Run binary directly with options:"
	@echo "  ./build/socks5-proxy -help   # Show all command-line options"
	@echo ""
	@echo "Examples:"
	@echo "  ./build/socks5-proxy -debug"
	@echo "  ./build/socks5-proxy -log-level debug"
	@echo "  ./build/socks5-proxy -port 8080"
	@echo "  ./build/socks5-proxy -config ./custom.properties"
.DEFAULT_GOAL := help
.PHONY: build install clean test lint lint-fix

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=pveidmapper

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/pveidmapper

# Install the application
install:
	$(GOCMD) install ./cmd/pveidmapper

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Update dependencies
update:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Run linter
lint:
	golangci-lint run

# Run linter with auto-fix
lint-fix:
	golangci-lint run --fix 

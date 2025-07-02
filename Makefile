.PHONY: build test clean help

build:
	go build -o bin/authenticator main.go

install:
	go build -o bin/authenticator main.go
	cp bin/authenticator /usr/local/bin/

# Clean build artifacts
clean:
	rm -rf bin/

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

# Cross-compile for different platforms
build-all: setup
	GOOS=linux GOARCH=amd64 go build -o bin/authenticator-linux main.go
	GOOS=windows GOARCH=amd64 go build -o bin/authenticator.exe main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/authenticator-mac main.go

# Show help
help:
	@echo "Available commands:"
	@echo "  build              - Build RTP Authenticator"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Install dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  build-all          - Cross-compile for all platforms"
	@echo "  help               - Show this help message"

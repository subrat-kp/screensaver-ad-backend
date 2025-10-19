.PHONY: swagger docs run test clean

# Generate Swagger documentation
swagger: docs

# Alias for swagger generation
docs:
	@echo "Generating Swagger documentation..."
	@export PATH=$(shell go env GOPATH)/bin:$$PATH && swag init --parseDependency --parseInternal
	@echo "✓ Swagger documentation generated successfully!"
	@echo "  - docs/swagger.json"
	@echo "  - docs/swagger.yaml"
	@echo "  - docs/docs.go"

# Run the application
run:
	@echo "Starting server..."
	@go run main.go


# Run the application in dev environment
run-dev:
	@echo "Starting development server.."
	@GO_ENV=dev go run main.go

# Run with automatic swagger generation
run-with-docs: docs run

# Run tests
test:
	@go test ./...

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf docs/
	@echo "✓ Cleaned!"

# Install swag tool if not already installed
install-swag:
	@echo "Installing swag CLI tool..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✓ Swag installed successfully!"

# Help command
help:
	@echo "Available commands:"
	@echo "  make swagger         - Generate Swagger documentation"
	@echo "  make docs            - Alias for 'make swagger'"
	@echo "  make run             - Run the application"
	@echo "  make run-with-docs   - Generate docs and run the application"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Remove generated documentation files"
	@echo "  make install-swag    - Install swag CLI tool"
	@echo "  make help            - Show this help message"

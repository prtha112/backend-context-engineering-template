# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go Clean Architecture template for building Product CRUD APIs with layered architecture. The project uses Gin for HTTP handling, PostgreSQL for data persistence, and follows Clean Architecture principles with clear separation of concerns.

## Development Commands

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/main.go

# Build the application
go build -o bin/app cmd/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/usecase -v

# Database migrations (using golang-migrate)
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run
```

## Architecture

This project follows **Clean Architecture** with these layers:

- **Domain** (`internal/domain/`): Core business entities, value objects, and domain services. No external dependencies.
- **Use Case** (`internal/usecase/`): Application business logic and orchestration. Depends only on domain interfaces.
- **Repository** (`internal/repository/`): Data access interfaces and implementations. Handles database operations.
- **Delivery** (`internal/delivery/`): HTTP handlers, DTOs, routing, and input validation using Gin framework.

### Dependency Flow
```
Delivery → Use Case → Domain ← Repository
```

### Key Dependencies
- **Gin**: HTTP web framework for REST API
- **lib/pq**: PostgreSQL driver
- **golang-migrate**: Database migration tool
- **validator/v10**: Request validation
- **logrus**: Structured logging
- **godotenv**: Environment configuration
- **OpenTelemetry**: Distributed tracing
- **testify**: Testing framework

## Configuration

Environment variables are loaded via `.env` file:
- `APP_NAME`, `APP_ENV`: Application metadata
- `HTTP_ADDR`, `HTTP_PORT`: Server configuration
- `DB_*`: Database connection parameters
- `LOG_LEVEL`: Logging configuration

## Testing Strategy

- **Unit tests**: Test domain entities and use cases with mock repositories
- **Integration tests**: Test repository implementations with real database
- Use `testify` for assertions and mocking
- Repository interfaces enable easy mocking for use case tests
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go Clean Architecture template for building Product CRUD APIs with layered architecture. The project uses Gin for HTTP handling, PostgreSQL for data persistence, and follows Clean Architecture principles with clear separation of concerns.

This template is specifically designed for context engineering and includes PRP (Project Requirement & Planning) documentation system for systematic feature development. The codebase is currently a template without actual implementation - use the PRPs and examples to guide feature development.

## Development Commands

```bash
# Install dependencies
go mod tidy

# Run the application (when implemented)
go run cmd/main.go

# Build the application
go build -o bin/app cmd/main.go

# Run tests
go test ./...

# Run tests with coverage and race detection
go test -cover -race ./...

# Run specific test package
go test ./internal/usecase -v

# Database migrations (using golang-migrate)
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down

# Code quality checks
go fmt ./...
go vet ./...
golangci-lint run
staticcheck ./...

# Full validation pipeline (as used in PRPs)
gofmt -l -w .
go mod tidy
golangci-lint run
go vet ./...
staticcheck ./...
go test ./... -v -race -count=1
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

## PRP (Project Requirement & Planning) System

This template uses a structured PRP system for feature development:

- **PRP Templates**: Located in `PRPs/templates/prp_base.md` - comprehensive template for feature planning
- **Generated PRPs**: Store feature-specific PRPs in `PRPs/` directory
- **Claude Commands**: Use `.claude/commands/generate-prp.md` to create new PRPs systematically

### PRP Development Process
1. Use the generate-prp command to analyze and create comprehensive feature PRPs
2. Follow the validation loops defined in PRPs (syntax, tests, integration)
3. Maintain strict Clean Architecture dependency rules
4. Include context-rich documentation and examples in PRPs

## Testing Strategy

- **Unit tests**: Test domain entities and use cases with mock repositories
- **Integration tests**: Test repository implementations with real database
- **Table-driven tests**: Follow Go testing patterns with structured test cases
- Use `testify` for assertions and mocking
- Repository interfaces enable easy mocking for use case tests
- Always run tests with `-race` flag to detect race conditions

## Implementation Guidelines

### Clean Architecture Rules
- **Domain** layer has no external dependencies
- **Use Case** layer depends only on domain interfaces
- **Repository** implementations are in infrastructure, interfaces in domain
- **Delivery** layer orchestrates HTTP/transport concerns

### Code Quality Requirements
- All code must pass: `golangci-lint run`, `go vet`, `staticcheck`
- Tests must pass with race detection: `go test -race ./...`
- Follow existing patterns when implementing new features
- Use structured logging with context propagation
- Implement graceful shutdown patterns
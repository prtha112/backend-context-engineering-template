# Backend Context Engineering Template

A complete Go Clean Architecture template for building scalable Product CRUD APIs with PostgreSQL, Docker deployment, and comprehensive testing. This project implements context engineering principles for systematic feature development.

## Overview

This template provides a production-ready foundation for Go backend services following Clean Architecture principles with clean separation of concerns, comprehensive testing, and deployment configurations.

**Current Status**: ✅ **Fully Implemented** - Complete Product CRUD API with PostgreSQL persistence, Docker deployment, and comprehensive test coverage.

## 🏗️ Architecture

This project follows **Clean Architecture** with these layers:

- **Domain** (`internal/domain/`): Core business entities, value objects, and domain services. No external dependencies.
- **Use Case** (`internal/usecase/`): Application business logic and orchestration. Depends only on domain interfaces.
- **Repository** (`internal/repository/`): Data access interfaces and implementations. Handles database operations.
- **Delivery** (`internal/delivery/`): HTTP handlers, DTOs, routing, and input validation using Gin framework.

### Dependency Flow
```
Delivery → Use Case → Domain ← Repository
```

## 🛠️ API Endpoints

- `POST /api/v1/products` - Create product with validation
- `GET /api/v1/products/:id` - Get single product by ID
- `GET /api/v1/products` - List products with pagination
- `PUT /api/v1/products/:id` - Update product with validation
- `DELETE /api/v1/products/:id` - Delete product by ID
- `GET /health` - Health check endpoint

## 🐳 Docker Deployment

### Quick Development Start
```bash
# Start PostgreSQL + pgAdmin, run migrations, start Go app
make dev-start
```

### Development (PostgreSQL only)
```bash
# Start PostgreSQL container
make dev-up

# Run migrations
make db-migrate-up

# Start Go app locally
go run cmd/main.go

# Stop containers
make dev-down
```

### Production Deployment
```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Database Management
```bash
# Run migrations
make db-migrate-up

# Rollback migrations
make db-migrate-down

# Force migration version
make db-migrate-force VERSION=1
```

## 🔧 Access Points

- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **pgAdmin**: http://localhost:5050 (admin@example.com / admin)
- **PostgreSQL**: localhost:5432

## 📊 Database Configuration

- **Database**: product_db
- **User**: app_user
- **Password**: app_password
- **Host**: localhost (development) / postgres (Docker)
- **Port**: 5432

## 🧪 Testing

### Run All Tests
```bash
# Unit tests with race detection
go test ./... -v -race -count=1

# Tests with coverage
go test -cover ./...

# Full validation pipeline
make validate
```

### Test Coverage
- **HTTP Handlers**: 90% coverage
- **Use Cases**: 54.9% coverage
- **Integration Tests**: Available for repository layer

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker setup)
- golang-migrate (for manual migrations)

### Local Development
1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd backend-context-engineering-template
   ```

2. **Start development environment**
   ```bash
   make dev-start
   ```

3. **Test the API**
   ```bash
   # Create a product
   curl -X POST http://localhost:8080/api/v1/products \
     -H "Content-Type: application/json" \
     -d '{"store_id": 1, "name": "Test Product", "amount": 10, "price": 29.99}'

   # Get all products
   curl http://localhost:8080/api/v1/products

   # Health check
   curl http://localhost:8080/health
   ```

### Production Deployment
1. **Configure environment variables** (copy from `.env.example`)
2. **Deploy with Docker Compose**
   ```bash
   docker-compose up -d
   ```

## 📁 Project Structure

```
/
├── cmd/
│   └── main.go                    # Application entry point
├── config/
│   └── config.go                  # Environment configuration
├── internal/
│   ├── domain/
│   │   ├── product.go             # Product entity with business rules
│   │   └── errors.go              # Domain-specific error types
│   ├── usecase/
│   │   ├── interfaces.go          # Repository interfaces (ports)
│   │   ├── product_usecase.go     # Business logic orchestration
│   │   └── product_usecase_test.go # Unit tests with mocks
│   ├── repository/
│   │   └── postgres/
│   │       ├── product_repository.go     # PostgreSQL implementation
│   │       └── product_repository_test.go # Integration tests
│   └── delivery/
│       └── http/
│           ├── dto/
│           │   └── product_dto.go         # Request/Response DTOs
│           ├── handlers/
│           │   ├── product_handler.go     # HTTP handlers
│           │   └── product_handler_test.go # Handler tests
│           ├── middleware/
│           │   ├── error_handler.go       # Global error handling
│           │   └── logger.go              # Request logging
│           └── router.go                  # Route definitions
├── migrations/
│   ├── 001_create_products_table.up.sql   # Database schema
│   └── 001_create_products_table.down.sql # Rollback script
├── pkg/
│   ├── database/
│   │   └── postgres.go            # Database connection setup
│   └── logger/
│       └── logger.go              # Structured logging setup
├── docker-compose.yaml            # Production deployment
├── docker-compose.dev.yaml        # Development environment
├── Dockerfile                     # Multi-stage build
├── Makefile                       # Build and deployment commands
└── .env                          # Environment variables
```

## ⚡ Key Features

### Clean Architecture Implementation
- **Strict dependency rules** maintained between layers
- **Interface-based** design for easy testing and mocking
- **Domain-driven** design with business rules in domain layer

### Production-Ready Features
- **Comprehensive input validation** with meaningful error messages
- **Structured logging** with request tracing
- **Graceful shutdown** handling
- **Database connection pooling** and health checks
- **SQL injection protection** with parameterized queries

### Development Experience
- **Comprehensive test suite** with mocks and integration tests
- **Docker development environment** with pgAdmin
- **Automated migrations** and database setup
- **Code quality tools** (gofmt, go vet, golangci-lint, staticcheck)
- **Make commands** for common tasks

### Security & Performance
- **Input validation** prevents invalid data entry
- **Parameterized queries** for SQL injection safety
- **Connection pooling** for database efficiency
- **Request timeouts** to prevent resource exhaustion
- **Structured error responses** without exposing internal errors

## 🔄 PRP Development System

This template uses a structured PRP (Project Requirement & Planning) system for feature development:

- **PRP Templates**: Located in `PRPs/templates/` - comprehensive templates for feature planning
- **Generated PRPs**: Store feature-specific PRPs in `PRPs/` directory
- **Context Engineering**: Systematic approach to feature development with validation loops

## 📋 Development Commands

```bash
# Install dependencies
go mod tidy

# Build the application
go build -o bin/app cmd/main.go

# Run the application
go run cmd/main.go

# Run tests
go test ./...

# Run tests with coverage and race detection
go test -cover -race ./...

# Code quality checks
go fmt ./...
go vet ./...
golangci-lint run
staticcheck ./...

# Full validation pipeline
make validate
```

## 🏆 Validation Results

- ✅ **Syntax & Style**: All code formatted and linted
- ✅ **Static Analysis**: Passes go vet and staticcheck
- ✅ **Unit Tests**: 90% handler coverage, all tests passing
- ✅ **Race Detection**: No race conditions detected
- ✅ **Build**: Application builds successfully
- ✅ **Integration**: API endpoints respond correctly

---

**Ready for production deployment with Docker or local development with comprehensive testing and validation!** 🎉


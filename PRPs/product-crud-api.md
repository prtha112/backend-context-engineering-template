# PRP: Product CRUD API with Go Clean Architecture

**Purpose**  
Implement a complete Product CRUD (Create, Read, Update, Delete) API using Go Clean Architecture principles with PostgreSQL persistence, comprehensive validation, logging, tracing, and testing.

---

## Core Principles
- **Context is King**: Include all relevant docs, references, and examples needed for implementation.  
- **Validation Loops**: Provide automated build/tests/linting to ensure correctness at each stage.  
- **Information Dense**: Use project-specific patterns, naming, and conventions.  
- **Progressive Success**: Implement in minimal working form first, then enhance.  
- **Follow Architecture Rules**: Maintain strict dependency boundaries between layers.

---

## Goal
**Build a production-ready Product CRUD API** following Clean Architecture principles with:
- RESTful endpoints for Product management (GET, POST, PUT, DELETE)
- Pagination, sorting, and filtering capabilities
- Input validation and error handling
- Structured logging and distributed tracing
- Database migrations and transaction management
- Comprehensive unit and integration tests
- Graceful shutdown and configuration management

---

## Why
- **Business Value**: Provides foundation for e-commerce product catalog management with scalable architecture
- **Integration**: Serves as template for other CRUD services in the system following established patterns  
- **Problem Solved**: Eliminates boilerplate setup and provides best-practice implementation for Go microservices

---

## What
- **User-visible behavior**: 
  - `GET /api/v1/products` - List products with pagination, sorting, filtering
  - `GET /api/v1/products/{id}` - Get single product by ID
  - `POST /api/v1/products` - Create new product
  - `PUT /api/v1/products/{id}` - Update existing product
  - `DELETE /api/v1/products/{id}` - Delete product
  - Structured JSON responses with consistent error handling
  - Request/response logging with correlation IDs

- **Technical requirements**: 
  - Clean Architecture with strict layer boundaries
  - PostgreSQL with connection pooling and transactions
  - Input validation using struct tags
  - OpenTelemetry distributed tracing
  - Structured logging with contextual fields
  - Database migrations for schema versioning
  - Environment-based configuration

---

## Success Criteria
- [ ] All CRUD endpoints respond correctly with proper HTTP status codes
- [ ] Request validation rejects invalid input with descriptive errors
- [ ] Database operations are transactional and handle concurrent access
- [ ] All tests pass: unit tests >90% coverage, integration tests for database layer
- [ ] No lint errors: `golangci-lint run` passes
- [ ] Performance: API responds <100ms for single-record operations
- [ ] Graceful shutdown completes within 30 seconds
- [ ] All layers maintain Clean Architecture dependency rules

---

## All Needed Context

**Documentation & References**
- `url:` https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html  
  `why:` Core Clean Architecture principles and dependency rules
- `url:` https://github.com/gin-gonic/gin  
  `why:` HTTP routing, middleware, JSON binding patterns
- `url:` https://pkg.go.dev/github.com/lib/pq  
  `why:` PostgreSQL connection management, error handling, transactions
- `url:` https://pkg.go.dev/github.com/golang-migrate/migrate/v4  
  `why:` Database migration patterns and CLI usage
- `url:` https://pkg.go.dev/github.com/go-playground/validator/v10  
  `why:` Struct validation tags and error handling
- `url:` https://pkg.go.dev/github.com/sirupsen/logrus  
  `why:` Structured logging with context fields
- `url:` https://opentelemetry.io/docs/instrumentation/go/  
  `why:` Distributed tracing and HTTP instrumentation

**Current Codebase Tree**  
```
/
├── CLAUDE.md                    # Project guidance and commands
├── go.mod                      # Dependencies already configured
├── examples/
│   └── README.md               # Clean Architecture structure guide
└── PRPs/
    └── templates/
        └── prp_base.md         # PRP template
```

**Desired Codebase Tree**  
```
/
├── cmd/
│   └── main.go                 # Application entry point and composition root
├── config/
│   ├── config.go               # Configuration struct and loading
│   └── database.go             # Database connection configuration
├── internal/
│   ├── domain/
│   │   ├── product.go          # Product entity with business rules
│   │   └── errors.go           # Domain-specific errors
│   ├── usecase/
│   │   ├── interfaces.go       # Repository and service interfaces
│   │   ├── product_usecase.go  # Business logic orchestration
│   │   └── product_usecase_test.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── product_repo.go # PostgreSQL implementation
│   │   │   └── product_repo_test.go
│   │   └── interfaces.go       # Repository interfaces
│   └── delivery/
│       └── http/
│           ├── handlers/
│           │   ├── product_handler.go # HTTP request handlers
│           │   └── product_handler_test.go
│           ├── middleware/
│           │   ├── logging.go  # Request logging middleware
│           │   └── tracing.go  # OpenTelemetry middleware
│           ├── dto/
│           │   └── product_dto.go # Request/response DTOs
│           └── router.go       # Route definitions
├── migrations/
│   ├── 001_create_products_table.up.sql
│   └── 001_create_products_table.down.sql
├── pkg/
│   ├── database/
│   │   └── postgres.go         # Database connection utilities
│   └── logger/
│       └── logger.go           # Logging setup and configuration
├── .env.example                # Environment variables template
├── Dockerfile                  # Container configuration
└── docker-compose.yml         # Development environment
```

**Known Gotchas**
- All repositories must implement interfaces defined in `internal/repository/interfaces.go` to enable dependency injection and testing
- Database transactions must be handled at the use case layer, not repository layer
- All HTTP handlers must use `*gin.Context` and follow consistent error response format
- Environment variables must be validated at startup to fail fast
- PostgreSQL parameterized queries use `$1, $2` format, not `?` placeholders
- OpenTelemetry middleware must be registered before other middleware to capture complete traces
- Logrus structured logging requires WithFields() for proper context
- Migration files require both .up.sql and .down.sql versions with timestamp prefixes

---

## Implementation Blueprint

### Data Models & Structure

**Domain Layer**: Pure business entities and rules
- Product entity with validation methods
- Domain-specific error types
- Business rule enforcement

**Use Case Layer**: Application business logic
- Product use case with CRUD operations
- Interface definitions for repositories
- Input/output ports for data transfer

**Repository Layer**: Data persistence abstraction  
- Repository interfaces in domain-facing package
- PostgreSQL implementation with connection pooling
- Transaction management and error translation

**Delivery Layer**: HTTP transport and presentation
- Gin handlers for REST endpoints
- Request/response DTOs with validation
- Middleware for logging and tracing

### Task List

**Task 1: Project Structure Setup**  
CREATE `cmd/main.go`
- Application entry point with dependency injection
- Graceful shutdown handling
- Configuration loading and validation

CREATE `config/config.go`  
- Environment variable mapping
- Configuration validation
- Database connection parameters

**Task 2: Domain Layer Implementation**  
CREATE `internal/domain/product.go`
- Product entity with fields: ID, Name, Description, Price, SKU, CreatedAt, UpdatedAt
- Validation methods for business rules
- Domain-specific methods (e.g., IsActive, SetPrice)

CREATE `internal/domain/errors.go`
- Domain error types (ErrProductNotFound, ErrInvalidPrice, etc.)
- Error wrapping utilities

**Task 3: Repository Layer**  
CREATE `internal/repository/interfaces.go`
- ProductRepository interface with CRUD methods
- Query parameters for filtering and pagination

CREATE `internal/repository/postgres/product_repo.go`
- PostgreSQL implementation of ProductRepository
- Connection management and query execution
- Error handling and transaction support

**Task 4: Use Case Layer**  
CREATE `internal/usecase/interfaces.go`
- Repository interface definitions
- External service interfaces

CREATE `internal/usecase/product_usecase.go`
- Business logic orchestration
- Transaction boundary management
- Error handling and logging

**Task 5: Delivery Layer**  
CREATE `internal/delivery/http/dto/product_dto.go`
- Request/response structures
- Validation tags and JSON mapping

CREATE `internal/delivery/http/handlers/product_handler.go`
- HTTP handlers for CRUD operations
- Request binding and validation
- Response formatting

CREATE `internal/delivery/http/router.go`
- Route registration and middleware setup
- CORS and security headers

**Task 6: Infrastructure Setup**  
CREATE `migrations/001_create_products_table.up.sql`
- Products table schema with indexes
- Constraints and data types

CREATE `pkg/database/postgres.go`
- Database connection setup
- Migration execution
- Health check utilities

CREATE `pkg/logger/logger.go`
- Structured logging configuration
- Context extraction utilities

**Task 7: Configuration & Environment**  
CREATE `.env.example`
- Environment variables template
- Development defaults

CREATE `Dockerfile`
- Multi-stage build for production
- Security best practices

**Task 8: Testing Implementation**  
CREATE test files for each layer:
- Domain unit tests
- Use case tests with mocked repositories  
- Repository integration tests with test database
- Handler tests with HTTP test utilities

---

## Integration Points

**DATABASE**  
- PostgreSQL with connection pooling via pgxpool
- Migration management using golang-migrate
- Transaction isolation levels for consistency
- Indexes on frequently queried columns (ID, SKU, CreatedAt)

**CONFIG**  
- Environment-based configuration with validation
- Graceful degradation for optional settings
- Secret management for database credentials

**ROUTES**  
```
GET    /api/v1/products           # List with pagination/filtering
GET    /api/v1/products/:id       # Get by ID
POST   /api/v1/products           # Create new
PUT    /api/v1/products/:id       # Update existing  
DELETE /api/v1/products/:id       # Delete by ID
GET    /health                    # Health check endpoint
```

**MIDDLEWARE STACK**
1. CORS headers
2. Request ID generation
3. OpenTelemetry tracing
4. Structured logging
5. Error recovery
6. Rate limiting (future enhancement)

---

## Validation Loop

**Level 1: Syntax & Style**  
```bash
go fmt ./...
go vet ./...
go mod tidy
golangci-lint run
```
Expected: No errors or warnings

**Level 2: Unit Tests**  
```bash
go test ./internal/domain/... -v
go test ./internal/usecase/... -v  
go test ./internal/delivery/... -v
```
Expected: All tests pass with >90% coverage

**Level 3: Integration Tests**  
```bash
go test ./internal/repository/... -v
```
Expected: Database integration tests pass

**Level 4: API Tests**  
```bash
go run cmd/main.go &
# Manual API testing with curl/Postman
curl -X GET http://localhost:8080/api/v1/products
curl -X POST http://localhost:8080/api/v1/products -d '{"name":"Test Product","price":99.99}'
```
Expected: All endpoints respond correctly

---

## Final Validation Checklist
- [ ] All tests pass: `go test ./...`
- [ ] No lint errors: `golangci-lint run`
- [ ] Database migrations apply and rollback correctly
- [ ] API endpoints handle happy path and error cases
- [ ] Request validation blocks invalid input
- [ ] Logs are structured with correlation IDs
- [ ] Traces appear in telemetry backend
- [ ] Graceful shutdown completes cleanly
- [ ] Configuration loads from environment variables
- [ ] Docker build completes successfully

---

## Anti-Patterns to Avoid
❌ Breaking dependency rules (delivery layer importing use case interfaces)  
❌ Database logic in handlers or use cases  
❌ Hardcoded configuration values  
❌ Missing transaction boundaries for data consistency  
❌ Generic error messages without context  
❌ Blocking operations without timeout contexts  
❌ Missing input validation on DTOs  
❌ SQL injection vulnerabilities  
❌ Resource leaks (unclosed connections, goroutines)  
❌ Logging sensitive data (passwords, tokens)

---

## Performance Considerations
- Connection pooling with appropriate limits
- Database query optimization with EXPLAIN ANALYZE
- Pagination to limit result set sizes
- Proper database indexes on query columns
- Context cancellation for request timeouts
- Structured logging to avoid string concatenation overhead

---

## Security Measures
- Input validation on all endpoints
- SQL injection prevention via parameterized queries
- HTTPS enforcement in production
- Sensitive data exclusion from logs
- Database credential protection
- CORS policy configuration
- Request size limits

---

**PRP Confidence Score: 9/10**

This PRP provides comprehensive context, specific implementation guidance, and executable validation steps. The success probability is high due to:
- Well-defined architecture boundaries
- Specific file structure and naming
- Concrete code examples and patterns
- Executable validation commands
- Detailed gotchas and anti-patterns
- External documentation references

The only minor risk is the complexity of coordinating all layers correctly, but the step-by-step approach and validation loops mitigate this concern.
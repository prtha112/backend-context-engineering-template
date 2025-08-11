# PRP – Product CRUD API Implementation (Go Clean Architecture)

**Purpose**  
Implement a complete Product CRUD API using Go Clean Architecture principles with Gin framework, PostgreSQL persistence, structured validation, comprehensive testing, and OpenTelemetry observability.

---

## Core Principles
- **Context is King**: Include all relevant docs, references, and examples needed for implementation.  
- **Validation Loops**: Provide automated build/tests/linting to ensure correctness at each stage.  
- **Information Dense**: Use project-specific patterns, naming, and conventions.  
- **Progressive Success**: Implement in minimal working form first, then enhance.  
- **Follow Architecture Rules**: Maintain strict dependency boundaries between layers.

---

## Goal
**Build a complete Product CRUD REST API with the following capabilities:**
- Create new products with validation
- Read products by ID and list all products with pagination
- Update existing products with partial/full updates
- Delete products with proper error handling
- Full Clean Architecture implementation with proper layer separation
- Comprehensive test coverage (unit + integration)
- Database migrations and schema management
- Structured logging and distributed tracing
- Production-ready error handling and validation

---

## Why
- **Business Value**: Provides foundational API for product management, enabling e-commerce operations, inventory tracking, and product catalog management.  
- **Integration**: Serves as the core template implementation that demonstrates Clean Architecture patterns for future features.  
- **Problem Solved**: Establishes consistent API patterns, testing approaches, and development practices for the entire application.

---

## What
- **User-visible behavior**: 
  - `POST /api/v1/products` - Create new product with JSON payload
  - `GET /api/v1/products/{id}` - Retrieve single product by ID
  - `GET /api/v1/products` - List products with optional pagination (?page=1&limit=10)
  - `PUT /api/v1/products/{id}` - Full product update
  - `PATCH /api/v1/products/{id}` - Partial product update
  - `DELETE /api/v1/products/{id}` - Delete product
- **Technical requirements**: 
  - JSON request/response format with consistent error handling
  - Request validation using struct tags and go-playground/validator
  - PostgreSQL persistence with proper indexing
  - Database migrations for schema management
  - Structured logging with request correlation IDs
  - OpenTelemetry tracing for observability
  - Graceful error responses with proper HTTP status codes
  - Unit and integration tests with >80% coverage

---

## Success Criteria
- All CRUD endpoints respond within <100ms for 95% of requests
- All validation errors return structured JSON with field-level details
- Database operations are properly transactional where needed
- All tests pass: `go test ./... -race -count=1`
- All linters pass: `golangci-lint run`, `go vet`, `staticcheck`
- Manual API testing demonstrates correct behavior for all endpoints
- All layers maintain Clean Architecture dependency rules
- Migrations can be applied and rolled back successfully
- Logs include proper correlation IDs and structured fields
- OpenTelemetry traces are generated for all HTTP requests

---

## All Needed Context

**Documentation & References**
- `url:` https://gin-gonic.com/docs/  
  `why:` HTTP routing, JSON handling, middleware patterns, request binding
- `url:` https://pkg.go.dev/github.com/lib/pq  
  `why:` PostgreSQL driver patterns, connection management
- `url:` https://github.com/golang-migrate/migrate  
  `why:` Database migration CLI usage and file structure patterns
- `url:` https://pkg.go.dev/github.com/go-playground/validator/v10  
  `why:` Struct tag validation, custom validators, error handling
- `url:` https://pkg.go.dev/github.com/stretchr/testify  
  `why:` Testing patterns, assertions, mocking, table-driven tests
- `url:` https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html  
  `why:` Clean Architecture principles, dependency rules, layer separation
- `url:` https://opentelemetry.io/docs/languages/go/getting-started/  
  `why:` OpenTelemetry instrumentation patterns for HTTP APIs
- `file:` /Users/sathabhronchangchuea/Desktop/Work/Freshket/backend-context-engineering-template/CLAUDE.md  
  `why:` Project conventions, directory structure, development commands, testing strategy
- `file:` /Users/sathabhronchangchuea/Desktop/Work/Freshket/backend-context-engineering-template/examples/README.md  
  `why:` Clean Architecture layer descriptions and folder structure patterns
- `file:` /Users/sathabhronchangchuea/Desktop/Work/Freshket/backend-context-engineering-template/go.mod  
  `why:` Dependency versions and module structure to follow

**Current Codebase Tree**  
```
/
├── CLAUDE.md                   # Project guidelines and commands
├── README.md                   # Project overview and setup
├── go.mod                      # Dependencies and module definition
├── INITIAL_product.md          # Feature requirements
├── PRPs/                       # Project requirement documents
│   └── templates/
│       └── prp_base.md        # PRP template
├── examples/                   # Architecture documentation
│   └── README.md              # Clean Architecture guide
└── .claude/                    # Claude configuration
    ├── settings.local.json
    └── commands/
        ├── generate-prp.md
        └── execute-prp.md
```

**Desired Codebase Tree**  
```
/
├── cmd/                        # Application entry point
│   └── main.go                # Main application with DI setup
├── config/                     # Configuration management
│   └── config.go              # Environment loading and validation
├── internal/                   # Core application layers
│   ├── domain/                # Domain entities (no dependencies)
│   │   ├── product.go         # Product entity with business rules
│   │   └── repository.go      # Repository interfaces
│   ├── usecase/               # Application logic (depends on domain)
│   │   ├── product_usecase.go # Product CRUD business logic
│   │   └── interfaces.go      # Use case interfaces/ports
│   ├── repository/            # Data access implementations
│   │   └── postgres/
│   │       ├── product_repo.go # PostgreSQL product repository
│   │       └── migrations.go   # Migration helper
│   └── delivery/              # HTTP transport layer
│       └── http/
│           ├── handlers/
│           │   └── product_handler.go # Gin HTTP handlers
│           ├── middleware/
│           │   ├── logging.go      # Request logging
│           │   ├── tracing.go      # OpenTelemetry middleware
│           │   └── cors.go         # CORS configuration
│           ├── dto/
│           │   └── product_dto.go  # Request/response DTOs
│           └── router.go           # Route configuration
├── migrations/                 # Database migration files
│   ├── 001_create_products_table.up.sql
│   └── 001_create_products_table.down.sql
├── pkg/                       # Shared packages
│   ├── logger/                # Structured logging setup
│   │   └── logger.go
│   ├── database/              # Database connection
│   │   └── postgres.go
│   └── tracing/               # OpenTelemetry setup
│       └── tracing.go
└── tests/                     # Integration tests
    └── integration/
        └── product_test.go
```

**Known Gotchas**
- All repository interfaces must be defined in `internal/domain/repository.go` to maintain dependency direction
- Database connections must be managed at the application level and injected into repositories
- All HTTP handlers must use context for request-scoped logging and tracing
- Validation errors must be transformed into consistent JSON error responses
- All database operations should use transactions where data consistency matters
- Migration files must follow golang-migrate naming convention: `{version}_{description}.{up|down}.sql`
- OpenTelemetry middleware must be registered before route handlers in Gin
- All tests must include cleanup for database state when using real DB connections
- Environment configuration must support both .env files and OS environment variables
- JSON struct tags must be consistent: `json:"field_name" validate:"required"`

---

## Implementation Blueprint

### Data Models & Structure
- **Domain Layer**: Product entity with business validation, repository interfaces
- **Use Case Layer**: CRUD operations with business logic, input/output DTOs
- **Repository Layer**: PostgreSQL implementation with proper error handling
- **Delivery Layer**: Gin HTTP handlers with request validation and response formatting

### Task List

**Task 1: Project Foundation**  
CREATE project structure and configuration
- Create directory structure following Clean Architecture
- Set up configuration management with environment variables
- Initialize logging and tracing infrastructure
- Create database connection management

**Task 2: Domain Layer**  
CREATE `internal/domain/product.go`  
- Define Product entity with ID, Name, Description, Price, CreatedAt, UpdatedAt
- Add domain validation methods (ValidateCreate, ValidateUpdate)
- Define business rules and constraints
- Create domain errors for validation failures

CREATE `internal/domain/repository.go`  
- Define ProductRepository interface with CRUD methods
- Include pagination parameters for List operations
- Define repository errors for not found, constraint violations

**Task 3: Database Layer**  
CREATE `migrations/001_create_products_table.up.sql`  
- Create products table with proper indexes
- Add constraints for business rules (price > 0, name not empty)
- Include created_at, updated_at timestamps

CREATE `migrations/001_create_products_table.down.sql`  
- Rollback script to drop products table

CREATE `pkg/database/postgres.go`  
- Database connection setup with connection pooling
- Health check functionality
- Transaction management utilities

CREATE `internal/repository/postgres/product_repo.go`  
- Implement ProductRepository interface
- Use lib/pq driver with database/sql
- Handle all CRUD operations with proper error wrapping
- Implement pagination for List operations
- Include SQL logging for debugging

**Task 4: Use Case Layer**  
CREATE `internal/usecase/interfaces.go`  
- Define use case interfaces and input/output ports
- Define request/response structs for each operation

CREATE `internal/usecase/product_usecase.go`  
- Implement business logic for CRUD operations
- Coordinate between domain validation and repository operations
- Handle error transformation from domain/repository to use case errors
- Include business logic like duplicate checking, soft deletes if needed

**Task 5: HTTP Delivery Layer**  
CREATE `internal/delivery/http/dto/product_dto.go`  
- Define request/response DTOs with validation tags
- Include proper JSON marshaling with snake_case
- Add omitempty tags for optional fields in responses

CREATE `internal/delivery/http/middleware/logging.go`  
- Request/response logging with correlation IDs
- Error logging with stack traces
- Performance metrics (request duration)

CREATE `internal/delivery/http/middleware/tracing.go`  
- OpenTelemetry integration with Gin
- Span creation and attribute setting
- Error recording in traces

CREATE `internal/delivery/http/handlers/product_handler.go`  
- Implement all CRUD HTTP handlers
- Use Gin context for request binding and validation
- Transform use case errors to appropriate HTTP status codes
- Include request/response logging

CREATE `internal/delivery/http/router.go`  
- Set up Gin router with middleware chain
- Define all product routes with proper HTTP methods
- Include health check endpoint
- Configure CORS if needed

**Task 6: Application Entry Point**  
CREATE `config/config.go`  
- Load environment variables with defaults
- Validate configuration on startup
- Support both .env files and OS environment variables

CREATE `cmd/main.go`  
- Application bootstrap with dependency injection
- Database connection setup and migration check
- HTTP server setup with graceful shutdown
- OpenTelemetry initialization and cleanup
- Error handling for startup failures

**Task 7: Comprehensive Testing**  
CREATE `internal/domain/product_test.go`  
- Unit tests for domain validation methods
- Test business rule enforcement
- Table-driven tests for edge cases

CREATE `internal/usecase/product_usecase_test.go`  
- Unit tests with mocked repository
- Test all CRUD operations
- Test error handling and validation
- Use testify for assertions and mocks

CREATE `internal/delivery/http/handlers/product_handler_test.go`  
- HTTP handler tests with test server
- Test request validation and error responses
- Test JSON marshaling/unmarshaling
- Use httptest for HTTP testing

CREATE `tests/integration/product_test.go`  
- End-to-end integration tests with real database
- Test complete request flow from HTTP to database
- Include transaction rollback for test isolation
- Test migration up/down operations

**Task 8: Documentation and Scripts**  
CREATE `.env.example`  
- Document all required environment variables
- Provide sensible defaults for development

UPDATE `README.md`  
- Add API endpoint documentation
- Include setup and running instructions
- Document testing approach

---

## Integration Points

**DATABASE**  
- PostgreSQL with connection pooling via lib/pq
- Migration management using golang-migrate CLI
- Proper indexing for query performance (products by ID, name)
- Transaction management for data consistency

**CONFIG**  
- Environment variables: DB_*, HTTP_*, LOG_LEVEL, OTEL_*
- Support for .env files in development
- Configuration validation on application startup

**ROUTES**  
- RESTful API design with consistent path patterns
- Proper HTTP status codes for different scenarios
- JSON request/response format with snake_case
- Error responses with detailed field validation

**OBSERVABILITY**  
- Structured logging with logrus and context correlation
- OpenTelemetry tracing for distributed observability
- Health check endpoints for monitoring
- Database connection health monitoring

---

## Validation Loop

**Level 1: Syntax & Style**  
```bash
gofmt -l -w .
go mod tidy
golangci-lint run
go vet ./...
staticcheck ./...
```
Expected: No errors/warnings, consistent formatting.

**Level 2: Unit Tests**  
```bash
go test ./internal/... -v -race -count=1
```
Expected: All tests pass, >80% coverage for business logic.

**Level 3: Integration Tests**  
```bash
go test ./tests/... -v -race -count=1
```
Expected: End-to-end functionality verified with real database.

**Level 4: Application Startup**  
```bash
go run cmd/main.go
```
Expected: Clean startup with no errors, database connection established.

**Level 5: API Manual Testing**  
```bash
curl -X POST localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","description":"Test Description","price":99.99}'

curl localhost:8080/api/v1/products/1
```
Expected: Proper JSON responses, error handling works.

---

## Final Validation Checklist
- [ ] All tests pass: `go test ./... -race -count=1`
- [ ] No lint errors: `golangci-lint run && go vet ./... && staticcheck ./...`
- [ ] Manual API calls work for all CRUD operations
- [ ] Database migrations apply and rollback successfully
- [ ] Error responses include proper HTTP status codes and JSON format
- [ ] Logs include correlation IDs and structured fields
- [ ] OpenTelemetry traces are visible for HTTP requests
- [ ] Application starts cleanly with environment configuration
- [ ] All layer dependencies follow Clean Architecture rules (domain → usecase → repository/delivery)
- [ ] Documentation updated with API examples and setup instructions

---

## Anti-Patterns to Avoid
❌ Breaking dependency rules (delivery/repository importing use cases or domain importing external libraries)  
❌ Skipping validation in request DTOs or domain entities  
❌ Hardcoding database connection strings instead of using environment config  
❌ Using panic/recover without proper logging and error context  
❌ Adding database queries directly in HTTP handlers (bypassing use cases)  
❌ Ignoring database transaction boundaries for multi-step operations  
❌ Missing correlation IDs in logs for request tracing  
❌ Inconsistent error response formats between endpoints  
❌ Using global variables for configuration or dependencies  
❌ Skipping migration rollback testing  

---

## Confidence Score: 9/10

This PRP provides comprehensive context, clear implementation tasks, executable validation steps, and follows established Go patterns. The structured approach with Clean Architecture ensures maintainable, testable code. All necessary external documentation is referenced, and the validation loops provide confidence in implementation quality.

The slight deduction is due to potential integration complexities with OpenTelemetry and the need for careful dependency injection setup, but these are well-documented challenges with clear solutions provided.
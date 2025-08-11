# PRP: Product CRUD API with Go Clean Architecture

**Purpose**  
Implement a complete Product CRUD API using Go Clean Architecture principles with Gin framework, PostgreSQL, comprehensive validation, logging, tracing, and testing.

---

## Core Principles
- **Context is King**: Include all relevant docs, references, and examples needed for implementation.  
- **Validation Loops**: Provide automated build/tests/linting to ensure correctness at each stage.  
- **Information Dense**: Use project-specific patterns, naming, and conventions.  
- **Progressive Success**: Implement in minimal working form first, then enhance.  
- **Follow Architecture Rules**: Maintain strict dependency boundaries between layers.

---

## Goal
Build a production-ready Product CRUD API with Clean Architecture that provides:
- Full CRUD operations (Create, Read, Update, Delete) for products
- Pagination, sorting, and filtering for product listings
- Request validation, structured logging, and distributed tracing
- Database migrations and multi-environment support
- Comprehensive unit and integration tests
- Graceful shutdown and configuration management

---

## Why
- **Business Value**: Core foundation for e-commerce/inventory systems requiring scalable product management
- **Integration**: Clean Architecture enables easy testing, maintenance, and future feature additions  
- **Problem Solved**: Provides maintainable, testable foundation for product data management with proper separation of concerns

---

## What
- **User-visible behavior**: 
  - `POST /products` - Create new product with validation
  - `GET /products` - List products with pagination, filtering, sorting
  - `GET /products/{id}` - Get single product by ID
  - `PUT /products/{id}` - Update existing product with validation
  - `DELETE /products/{id}` - Delete product by ID
- **Technical requirements**: JSON API, PostgreSQL storage, request validation, structured logging, OpenTelemetry tracing

---

## Success Criteria
- All CRUD endpoints respond correctly with proper HTTP status codes
- Pagination works with `limit`, `offset`, `sort` query parameters
- Validation errors return structured error responses
- All requests are traced with OpenTelemetry
- Structured logs include request context and performance metrics
- Database migrations run successfully up and down
- Unit tests achieve >80% coverage with mocked dependencies
- Integration tests verify end-to-end functionality
- Graceful shutdown handles SIGTERM/SIGINT properly

---

## All Needed Context

**Documentation & References**
- `url:` https://pkg.go.dev/github.com/gin-gonic/gin  
  `why:` HTTP routing, middleware patterns, context handling
- `url:` https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md  
  `why:` PostgreSQL migration patterns and connection strings
- `url:` https://signoz.io/blog/opentelemetry-gin/  
  `why:` OpenTelemetry Gin instrumentation examples
- `url:` https://pkg.go.dev/github.com/go-playground/validator/v10  
  `why:` Struct validation tags and patterns
- `url:` https://github.com/sirupsen/logrus  
  `why:` Structured logging with fields and context
- `url:` https://github.com/stretchr/testify  
  `why:` Assertions, mocks, and testing patterns
- `url:` https://github.com/bxcodec/go-clean-arch  
  `why:` Clean Architecture reference implementation
- `url:` https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1  
  `why:` Go-specific Clean Architecture patterns

**Current Codebase Tree**  
```
/Users/sathabhronchangchuea/Desktop/Work/Freshket/backend-context-engineering-template/
├── CLAUDE.md
├── go.mod (with dependencies: gin, lib/pq, golang-migrate, validator, logrus, testify, otel)
├── examples/README.md (Clean Architecture folder structure reference)
└── PRPs/templates/prp_base.md
```

**Desired Codebase Tree**  
```
/Users/sathabhronchangchuea/Desktop/Work/Freshket/backend-context-engineering-template/
├── cmd/
│   └── main.go                              # Application entry point with DI setup
├── config/
│   └── config.go                           # Environment configuration loading
├── internal/
│   ├── domain/
│   │   ├── product.go                      # Product entity with business rules
│   │   └── errors.go                       # Domain-specific error definitions
│   ├── usecase/
│   │   ├── interfaces.go                   # UseCase and Repository interfaces
│   │   └── product_usecase.go              # Product business logic orchestration
│   ├── repository/
│   │   └── postgres/
│   │       └── product_repository.go       # PostgreSQL implementation
│   └── delivery/
│       └── http/
│           ├── dto/
│           │   └── product_dto.go          # Request/Response DTOs
│           ├── handler/
│           │   ├── product_handler.go      # HTTP handlers
│           │   └── response.go             # Standard response wrapper
│           ├── middleware/
│           │   ├── logger.go              # Request logging middleware
│           │   └── tracer.go              # OpenTelemetry middleware
│           └── router.go                   # Route definitions
├── migrations/
│   ├── 000001_create_products.up.sql      # Create products table
│   └── 000001_create_products.down.sql    # Drop products table
├── .env.example                            # Environment template
└── go.mod                                  # Dependencies already defined
```

**Known Gotchas**
- Clean Architecture: Inner layers (domain) cannot import outer layers (delivery/infrastructure)
- All repositories must implement interfaces defined in usecase layer
- Use context.Context for request-scoped data (logging, tracing)
- Database transactions should be managed at UseCase level, not Repository level
- JSON struct tags must be consistent with validation tags
- OpenTelemetry middleware must be registered before route handlers
- Environment variables need validation and defaults
- Graceful shutdown requires proper signal handling and cleanup

---

## Implementation Blueprint

### Data Models & Structure

**Domain Layer Pattern** (from research):
```go
// Product entity with business rules and validation
type Product struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name" validate:"required,min=3,max=100"`
    Description string    `json:"description" validate:"max=500"`
    Price       float64   `json:"price" validate:"required,gt=0"`
    SKU         string    `json:"sku" validate:"required,unique"`
    Stock       int       `json:"stock" validate:"gte=0"`
    CategoryID  uint      `json:"category_id" validate:"required"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**UseCase Layer Pattern** (from research):
```go
type ProductUseCase interface {
    CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
    GetProduct(ctx context.Context, id uint) (*Product, error)
    UpdateProduct(ctx context.Context, id uint, req *UpdateProductRequest) (*Product, error)
    DeleteProduct(ctx context.Context, id uint) error
    ListProducts(ctx context.Context, params ListProductsParams) (*ProductList, error)
}

type ProductRepository interface {
    Create(ctx context.Context, product *Product) (*Product, error)
    GetByID(ctx context.Context, id uint) (*Product, error)
    Update(ctx context.Context, product *Product) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, params ListParams) ([]*Product, int64, error)
}
```

**Repository Layer Pattern** (from PostgreSQL research):
```go
func (r *productRepository) Create(ctx context.Context, product *Product) (*Product, error) {
    query := `INSERT INTO products (name, description, price, sku, stock, category_id) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowContext(ctx, query, product.Name, product.Description, 
        product.Price, product.SKU, product.Stock, product.CategoryID).
        Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
    return product, err
}
```

**Delivery Layer Pattern** (from Gin research):
```go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    product, err := h.productUC.CreateProduct(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, SuccessResponse{Data: product})
}
```

### Task List

**Task 1: Project Structure & Configuration**  
CREATE project directory structure and configuration
- `cmd/main.go` - Application entry point with DI
- `config/config.go` - Environment config with validation
- `.env.example` - Environment template

**Task 2: Domain Layer Implementation**  
CREATE `internal/domain/product.go` and `internal/domain/errors.go`
- Product entity with validation tags
- Domain-specific error types
- Business rule validation methods

**Task 3: Interface Definitions**  
CREATE `internal/usecase/interfaces.go`  
- ProductUseCase interface with all CRUD methods
- ProductRepository interface for data persistence
- Request/Response struct definitions

**Task 4: UseCase Layer Implementation**  
CREATE `internal/usecase/product_usecase.go`
- Implement ProductUseCase interface
- Business logic orchestration with error handling
- Transaction management and validation

**Task 5: Repository Layer Implementation**  
CREATE `internal/repository/postgres/product_repository.go`
- Implement ProductRepository interface
- PostgreSQL queries with proper error handling  
- Context-aware database operations

**Task 6: Database Migrations**  
CREATE migration files in `migrations/`
- `000001_create_products.up.sql` with products table schema
- `000001_create_products.down.sql` with table drop
- Include indexes for performance

**Task 7: HTTP Delivery Layer - DTOs**  
CREATE `internal/delivery/http/dto/product_dto.go`
- Request DTOs with validation tags
- Response DTOs for API consistency
- Pagination and filtering structs

**Task 8: HTTP Delivery Layer - Handlers**  
CREATE `internal/delivery/http/handler/product_handler.go` and `response.go`
- CRUD endpoint handlers with Gin context
- Request validation and error handling
- Standard response format implementation

**Task 9: Middleware Implementation**  
CREATE middleware in `internal/delivery/http/middleware/`
- `logger.go` - Structured logging with Logrus
- `tracer.go` - OpenTelemetry instrumentation
- Context propagation for request tracking

**Task 10: HTTP Router Setup**  
CREATE `internal/delivery/http/router.go`
- Route definitions with proper HTTP methods
- Middleware registration order
- Route grouping and versioning

**Task 11: Application Wiring**  
MODIFY `cmd/main.go` to wire all components
- Database connection and migration
- Dependency injection setup
- HTTP server with graceful shutdown
- Configuration loading and validation

**Task 12: Comprehensive Testing**  
CREATE test files for each layer
- `internal/domain/product_test.go` - Entity validation tests
- `internal/usecase/product_usecase_test.go` - Business logic tests with mocks
- `internal/delivery/http/handler/product_handler_test.go` - HTTP handler tests
- Use testify for assertions and mocking

---

## Integration Points

**DATABASE**  
- PostgreSQL with connection pooling
- Migration management with golang-migrate
- Transaction support for complex operations
- Proper indexes for query performance

**CONFIG**  
- Environment-based configuration with godotenv
- Validation of required environment variables
- Default values for development environment
- Support for different environments (dev/staging/prod)

**ROUTES**  
- RESTful endpoint design following HTTP standards
- Proper HTTP status codes for different scenarios
- Content-Type and Accept header handling
- CORS support for frontend integration

**OBSERVABILITY**
- OpenTelemetry tracing with request correlation
- Structured logging with contextual information
- Health check endpoints for monitoring
- Metrics collection for performance monitoring

---

## Validation Loop

**Level 1: Syntax & Style**  
```bash
go mod tidy
go fmt ./...
go vet ./...
# Optional: golangci-lint run (if available)
```
Expected: No errors/warnings, dependencies resolved.

**Level 2: Build & Compilation**  
```bash
go build -o bin/app cmd/main.go
```
Expected: Clean compilation with executable created.

**Level 3: Unit Tests**  
```bash
go test ./... -v -race -count=1
go test -cover ./...
```
Expected: All tests pass, >80% coverage achieved.

**Level 4: Integration Tests**  
```bash
# Start test database
# Run migrations
go test ./internal/delivery/http/handler -v -tags=integration
```
Expected: End-to-end functionality verified.

**Level 5: Manual API Testing**  
```bash
# Start application
go run cmd/main.go

# Test endpoints
curl -X POST http://localhost:8080/products -d '{"name":"Test Product","price":99.99,"sku":"TEST001","stock":10,"category_id":1}' -H "Content-Type: application/json"
curl http://localhost:8080/products
curl http://localhost:8080/products/1
```
Expected: Proper JSON responses with correct status codes.

---

## Final Validation Checklist
- [ ] All CRUD endpoints working with proper HTTP methods
- [ ] Request validation returns structured error messages  
- [ ] Pagination, sorting, filtering work on GET /products
- [ ] All tests pass: `go test ./...`
- [ ] No lint errors: `go vet ./...` and `golangci-lint run`
- [ ] Database migrations run up and down successfully
- [ ] Structured logs include request IDs and contextual information
- [ ] OpenTelemetry traces are generated for all requests
- [ ] Graceful shutdown works properly on SIGTERM/SIGINT
- [ ] Configuration loads from environment variables
- [ ] Clean Architecture boundaries maintained (no circular imports)

---

## Anti-Patterns to Avoid
❌ Breaking Clean Architecture dependency rules (outer importing inner)  
❌ Skipping request validation at HTTP boundary  
❌ Using domain entities directly in HTTP responses  
❌ Managing database transactions in repository layer  
❌ Hardcoding configuration values instead of environment variables  
❌ Ignoring context cancellation in long-running operations  
❌ Missing error handling and logging for debugging  
❌ Using panic/recover without proper logging  
❌ Creating tight coupling between layers  
❌ Missing indexes on frequently queried database columns

---

## Context-Rich Implementation Examples

**Gin Handler Pattern** (from research):
```go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // Extract request ID for logging
    requestID := c.GetHeader("X-Request-ID")
    if requestID == "" {
        requestID = uuid.New().String()
    }
    
    // Add to context for logging
    ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
    
    // Structured logging with context
    logger := logrus.WithFields(logrus.Fields{
        "request_id": requestID,
        "method":     c.Request.Method,
        "path":       c.Request.URL.Path,
    })
    
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.WithError(err).Error("Invalid request payload")
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request payload",
            Details: getValidationErrors(err),
        })
        return
    }
    
    product, err := h.productUC.CreateProduct(ctx, &req)
    if err != nil {
        logger.WithError(err).Error("Failed to create product")
        c.JSON(getHTTPStatus(err), ErrorResponse{Error: err.Error()})
        return
    }
    
    logger.WithField("product_id", product.ID).Info("Product created successfully")
    c.JSON(http.StatusCreated, SuccessResponse{Data: product})
}
```

**OpenTelemetry Setup Pattern** (from research):
```go
// In main.go or middleware setup
func setupTracing() {
    r := gin.New()
    r.Use(otelgin.Middleware("product-crud-api"))
    
    // Custom span attributes
    r.Use(func(c *gin.Context) {
        span := trace.SpanFromContext(c.Request.Context())
        span.SetAttributes(
            attribute.String("http.user_agent", c.GetHeader("User-Agent")),
            attribute.String("request.id", c.GetHeader("X-Request-ID")),
        )
        c.Next()
    })
}
```

**Validation Pattern** (from research):
```go
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    SKU         string  `json:"sku" validate:"required,alphanum,min=3,max=50"`
    Stock       int     `json:"stock" validate:"gte=0"`
    CategoryID  uint    `json:"category_id" validate:"required,gt=0"`
}

func getValidationErrors(err error) map[string]string {
    errors := make(map[string]string)
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            errors[e.Field()] = getValidationMessage(e)
        }
    }
    return errors
}
```

**Database Migration Pattern** (from research):
```sql
-- 000001_create_products.up.sql
BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    sku VARCHAR(50) NOT NULL UNIQUE,
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);

COMMIT;
```

**Testify Testing Pattern** (from research):
```go
func TestProductUseCase_CreateProduct(t *testing.T) {
    tests := []struct {
        name          string
        request       *CreateProductRequest
        mockSetup     func(*mocks.ProductRepository)
        expectedError error
    }{
        {
            name: "successful product creation",
            request: &CreateProductRequest{
                Name:       "Test Product",
                Price:      99.99,
                SKU:        "TEST001",
                Stock:      10,
                CategoryID: 1,
            },
            mockSetup: func(mockRepo *mocks.ProductRepository) {
                mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).
                    Return(&domain.Product{ID: 1, Name: "Test Product"}, nil)
            },
            expectedError: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(mocks.ProductRepository)
            tt.mockSetup(mockRepo)
            
            uc := NewProductUseCase(mockRepo)
            result, err := uc.CreateProduct(context.Background(), tt.request)
            
            if tt.expectedError != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedError, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
            }
            
            mockRepo.AssertExpectations(t)
        })
    }
}
```

---

## Confidence Score: 9/10

This PRP provides comprehensive context for one-pass implementation success with:
- ✅ Complete file structure and implementation patterns
- ✅ Concrete code examples from extensive library research
- ✅ Sequential task breakdown following Clean Architecture principles
- ✅ Executable validation gates for each development stage
- ✅ Rich context about gotchas, conventions, and best practices
- ✅ Integration examples for all major dependencies (Gin, PostgreSQL, OpenTelemetry, etc.)

The high confidence score reflects thorough research and detailed implementation guidance that should enable successful development without additional context requirements.
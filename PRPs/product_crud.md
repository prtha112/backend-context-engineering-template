# PRP - Product CRUD Implementation

**Purpose**  
Implement a complete Product CRUD API following Clean Architecture principles with Gin framework, PostgreSQL persistence, and comprehensive validation. This PRP provides all necessary context for one-pass implementation success.

---

## Core Principles
- **Context is King**: Include all relevant docs, references, and examples needed for implementation.  
- **Validation Loops**: Provide automated build/tests/linting to ensure correctness at each stage.  
- **Information Dense**: Use project-specific patterns, naming, and conventions.  
- **Progressive Success**: Implement in minimal working form first, then enhance.  
- **Follow Architecture Rules**: Maintain strict dependency boundaries between layers.

---

## Goal
**Implement a complete Product CRUD API with PostgreSQL persistence following Clean Architecture**

Create endpoints for Create, Read (by ID and list), Update, and Delete operations for Product entities. The system must handle the specified database schema with proper validation, error handling, and follow Clean Architecture dependency rules.

---

## Why
- **Business Value**: Enables product management functionality as the foundation for an e-commerce system.  
- **Integration**: Serves as the template implementation for other CRUD operations in the system.  
- **Problem Solved**: Provides structured, maintainable product data management with proper separation of concerns.

---

## What
**User-visible behavior:**  
- `POST /products` - Create new product with validation
- `GET /products/:id` - Retrieve single product by ID  
- `GET /products` - List all products with optional pagination
- `PUT /products/:id` - Update existing product with validation
- `DELETE /products/:id` - Delete product by ID
- Proper HTTP status codes and JSON error responses
- Input validation with meaningful error messages

**Technical requirements:**  
- PostgreSQL table: `products` with fields (id, store_id, name, description, amount, price)
- Clean Architecture layers: Domain → Use Case → Repository → Delivery
- JSON request/response DTOs with validation tags
- Database migrations for schema creation
- Unit and integration tests
- Graceful error handling and structured logging

---

## Success Criteria
- All CRUD endpoints respond correctly with proper HTTP status codes
- Input validation prevents invalid data entry
- Database operations use parameterized queries (SQL injection safe)  
- All tests pass: `go test ./... -v -race -count=1`
- All linting passes: `gofmt`, `go vet`, `golangci-lint run`, `staticcheck`
- Clean Architecture dependency rules maintained
- API response time <100ms for single operations (excluding network latency)

---

## All Needed Context

**Documentation & References**
- `url:` https://go.dev/doc/effective_go  
  `why:` Error handling patterns, testing best practices, project structure
- `url:` https://pkg.go.dev/github.com/gin-gonic/gin  
  `why:` HTTP handlers, middleware, JSON binding patterns  
- `url:` https://pkg.go.dev/github.com/lib/pq  
  `why:` PostgreSQL connection handling, parameterized queries, error types
- `url:` https://pkg.go.dev/github.com/go-playground/validator/v10  
  `why:` Struct validation tags and error handling
- `file:` `/examples/README.md`  
  `why:` Clean Architecture folder structure and layer responsibilities
- `file:` `/CLAUDE.md`  
  `why:` Development commands, architecture rules, dependencies

**Current Codebase Tree**  
```
/
├── CLAUDE.md
├── examples/
│   └── README.md
├── PRPs/
│   ├── templates/
│   │   └── prp_base.md
│   └── product_crud.md
└── INITIAL_product.md
```

**Desired Codebase Tree**  
```
/
├── cmd/
│   └── main.go                    # Application entry point with DI setup
├── config/
│   └── config.go                  # Environment variable loading
├── internal/
│   ├── domain/
│   │   ├── product.go             # Product entity with business rules
│   │   └── errors.go              # Domain-specific error types
│   ├── usecase/
│   │   ├── interfaces.go          # Repository interfaces (ports)
│   │   ├── product_usecase.go     # Business logic orchestration
│   │   └── product_usecase_test.go # Unit tests with mocks
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── product_repository.go     # PostgreSQL implementation
│   │   │   └── product_repository_test.go # Integration tests
│   │   └── interfaces.go          # Repository interfaces
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
│   └── 001_create_products_table.up.sql   # Database schema
├── pkg/
│   ├── database/
│   │   └── postgres.go            # Database connection setup
│   └── logger/
│       └── logger.go              # Structured logging setup
├── go.mod                         # Dependencies
├── go.sum
├── .env.example                   # Environment template
└── Makefile                       # Build and test commands
```

**Known Gotchas**
- `store_id` must NOT have foreign key constraint - just regular integer field
- All repository interfaces must be defined in the use case layer, not repository layer
- Use parameterized queries ($1, $2) with lib/pq, never string concatenation
- Gin context timeout should be propagated to database operations
- Use `database/sql.NullString`, `sql.NullInt64` for nullable fields
- PostgreSQL numeric type requires `database/sql/driver.Valuer` interface for decimal handling
- All HTTP handlers must validate input using struct tags before calling use cases
- Error responses should be consistent JSON format across all endpoints

---

## Implementation Blueprint

### Database Schema & Migration
```sql
-- migrations/001_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    amount INTEGER NOT NULL DEFAULT 0,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_store_id ON products(store_id);

-- migrations/001_create_products_table.down.sql  
DROP TABLE IF EXISTS products;
```

### Data Models & Structure

**Domain Layer** (`internal/domain/`):
- `Product` entity with business validation methods
- Domain-specific errors (ProductNotFound, InvalidProduct)
- Value objects if needed (Price, ProductName)

**Use Case Layer** (`internal/usecase/`):
- Repository interfaces (ProductRepository)  
- ProductUseCase struct with CRUD methods
- Input/output DTOs for use case boundaries
- Business logic orchestration

**Repository Layer** (`internal/repository/postgres/`):
- ProductRepository implementation
- Database connection handling
- SQL query implementations with proper error handling

**Delivery Layer** (`internal/delivery/http/`):
- HTTP handlers for each CRUD operation
- Request/Response DTOs with validation tags
- Route definitions and middleware setup
- Error response formatting

### Task List

**Task 1: Project Structure Setup**  
CREATE directory structure and basic files:
- `cmd/main.go` - Application entry point
- `config/config.go` - Environment configuration
- `pkg/database/postgres.go` - Database connection
- `pkg/logger/logger.go` - Structured logging
- `go.mod` - Initialize Go module with dependencies

**Task 2: Domain Layer**  
CREATE `internal/domain/product.go`:
```go
type Product struct {
    ID          int64           `json:"id" db:"id"`
    StoreID     int64           `json:"store_id" db:"store_id"`
    Name        string          `json:"name" db:"name"`
    Description sql.NullString  `json:"description" db:"description"`
    Amount      int64           `json:"amount" db:"amount"`
    Price       float64         `json:"price" db:"price"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// Business validation methods
func (p *Product) Validate() error
func (p *Product) IsValidPrice() bool
```

CREATE `internal/domain/errors.go`:
- Define domain-specific errors (ErrProductNotFound, ErrInvalidProduct)

**Task 3: Use Case Layer**  
CREATE `internal/usecase/interfaces.go`:
```go
type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
    GetByID(ctx context.Context, id int64) (*domain.Product, error)
    GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error)
    Update(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error)
    Delete(ctx context.Context, id int64) error
}
```

CREATE `internal/usecase/product_usecase.go`:
- ProductUseCase struct with repository dependency
- Implement all CRUD operations with business logic
- Handle errors and logging

**Task 4: Repository Implementation**  
CREATE `internal/repository/postgres/product_repository.go`:
- Implement ProductRepository interface
- Use parameterized queries with lib/pq
- Handle PostgreSQL-specific errors
- Include database transaction handling for complex operations

**Task 5: HTTP Delivery Layer**  
CREATE `internal/delivery/http/dto/product_dto.go`:
```go
type CreateProductRequest struct {
    StoreID     int64   `json:"store_id" binding:"required,min=1"`
    Name        string  `json:"name" binding:"required,min=1,max=100"`
    Description string  `json:"description" binding:"max=1000"`
    Amount      int64   `json:"amount" binding:"required,min=0"`
    Price       float64 `json:"price" binding:"required,min=0"`
}

type ProductResponse struct {
    ID          int64  `json:"id"`
    StoreID     int64  `json:"store_id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Amount      int64  `json:"amount"`
    Price       float64 `json:"price"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}
```

CREATE `internal/delivery/http/handlers/product_handler.go`:
```go
type ProductHandler struct {
    usecase usecase.ProductUseCase
    logger  *logrus.Logger
}

func (h *ProductHandler) CreateProduct(c *gin.Context)
func (h *ProductHandler) GetProduct(c *gin.Context)
func (h *ProductHandler) GetProducts(c *gin.Context)
func (h *ProductHandler) UpdateProduct(c *gin.Context)
func (h *ProductHandler) DeleteProduct(c *gin.Context)
```

**Task 6: Database Migration**  
CREATE `migrations/001_create_products_table.up.sql` and `.down.sql`
- Include proper indexes for performance
- Handle PostgreSQL-specific data types

**Task 7: Application Wiring**  
MODIFY `cmd/main.go`:
- Database connection setup
- Dependency injection for all layers  
- HTTP server configuration with Gin
- Graceful shutdown handling

**Task 8: Router Configuration**  
CREATE `internal/delivery/http/router.go`:
```go
func SetupRouter(productHandler *handlers.ProductHandler) *gin.Engine {
    r := gin.New()
    r.Use(middleware.Logger())
    r.Use(middleware.ErrorHandler())
    
    api := r.Group("/api/v1")
    {
        products := api.Group("/products")
        {
            products.POST("", productHandler.CreateProduct)
            products.GET("/:id", productHandler.GetProduct)
            products.GET("", productHandler.GetProducts)
            products.PUT("/:id", productHandler.UpdateProduct)
            products.DELETE("/:id", productHandler.DeleteProduct)
        }
    }
    return r
}
```

---

## Integration Points

**DATABASE**  
- PostgreSQL connection via lib/pq driver
- Migration management using golang-migrate
- Connection pooling via database/sql
- Transaction handling for atomic operations

**CONFIG**  
Environment variables (`.env`):
```env
APP_NAME=product-service
APP_ENV=development
HTTP_ADDR=0.0.0.0
HTTP_PORT=8080

DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=product_db
DB_SSLMODE=disable

LOG_LEVEL=info
```

**ROUTES**  
- RESTful endpoint design: `/api/v1/products`
- Consistent HTTP methods (POST, GET, PUT, DELETE)
- Proper HTTP status codes (200, 201, 400, 404, 500)

**DEPENDENCIES**  
Required Go modules:
```go
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/lib/pq v1.10.9
    github.com/golang-migrate/migrate/v4 v4.16.2
    github.com/go-playground/validator/v10 v10.15.5
    github.com/sirupsen/logrus v1.9.3
    github.com/joho/godotenv v1.4.0
    github.com/stretchr/testify v1.8.4
)
```

---

## Validation Loop

**Level 1: Syntax & Style**  
```bash
# Format code
gofmt -l -w .

# Clean up dependencies  
go mod tidy

# Static analysis
golangci-lint run
go vet ./...
staticcheck ./...
```
Expected: No errors/warnings.

**Level 2: Unit Tests**  
```bash
# Run all tests with race detection
go test ./... -v -race -count=1

# Test coverage
go test -cover ./...
```

Create test files:
- `internal/usecase/product_usecase_test.go` - Mock repository tests
- `internal/delivery/http/handlers/product_handler_test.go` - HTTP handler tests  
- `internal/repository/postgres/product_repository_test.go` - Integration tests

**Level 3: Integration Tests**
```bash
# API testing with curl/httpie
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"store_id": 1, "name": "Test Product", "amount": 10, "price": 29.99}'
```

---

## Testing Patterns

**Unit Tests with Mocks:**
```go
func TestProductUseCase_CreateProduct(t *testing.T) {
    tests := []struct {
        name    string
        product *domain.Product
        mockFn  func(*mocks.ProductRepository)
        want    *domain.Product
        wantErr bool
    }{
        {
            name: "successful creation",
            product: &domain.Product{Name: "Test Product", Price: 29.99},
            mockFn: func(m *mocks.ProductRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(&domain.Product{ID: 1}, nil)
            },
            want: &domain.Product{ID: 1},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &mocks.ProductRepository{}
            tt.mockFn(repo)
            
            uc := usecase.NewProductUseCase(repo, logger)
            got, err := uc.CreateProduct(context.Background(), tt.product)
            
            assert.Equal(t, tt.wantErr, err != nil)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**Integration Tests:**
```go  
func TestProductRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewProductRepository(db, logger)
    
    // Test Create
    product := &domain.Product{Name: "Integration Test Product", Price: 19.99}
    created, err := repo.Create(context.Background(), product)
    
    require.NoError(t, err)
    assert.NotZero(t, created.ID)
    assert.Equal(t, product.Name, created.Name)
}
```

---

## Final Validation Checklist
- [ ] All tests pass: `go test ./... -v -race -count=1`
- [ ] No lint errors: `golangci-lint run && go vet ./... && staticcheck ./...`
- [ ] Manual API calls work as expected (all CRUD operations)
- [ ] Errors handled gracefully with proper HTTP status codes
- [ ] Logs are structured and contextual with request tracing
- [ ] Database queries use parameterized statements
- [ ] Clean Architecture dependency rules maintained
- [ ] Input validation prevents invalid data with meaningful error messages
- [ ] Database migrations run successfully up and down

---

## Anti-Patterns to Avoid
❌ **Breaking dependency rules** - Domain layer importing delivery/repository  
❌ **String concatenation in SQL queries** - Use parameterized queries ($1, $2)  
❌ **Ignoring database errors** - Always handle sql.ErrNoRows appropriately  
❌ **Hardcoding database config** - Use environment variables  
❌ **Catch-all panic recovery** - Log panics with context before recovery  
❌ **Exposing domain entities in HTTP layer** - Use DTOs for API boundaries  
❌ **Missing input validation** - Validate all incoming data with struct tags  
❌ **Inconsistent error response format** - Standardize error JSON structure  
❌ **Missing database indexes** - Index foreign keys and frequently queried fields  
❌ **Ignoring context cancellation** - Propagate context through all layers  

---

## Performance Considerations
- Use database connection pooling (default with database/sql)
- Add indexes on store_id and other frequently queried fields  
- Implement pagination for list endpoints to prevent large result sets
- Consider adding database query timeout using context.WithTimeout
- Use EXPLAIN ANALYZE to validate query performance in PostgreSQL

---

## Security Checklist
- All database queries use parameterized statements (SQL injection prevention)
- Input validation on all incoming data using struct tags  
- Proper HTTP status codes (don't expose internal errors)
- Environment variables for sensitive configuration (database credentials)
- Request timeout handling to prevent resource exhaustion
- Structured logging without exposing sensitive data

---

**PRP Confidence Level: 9/10**

This PRP provides comprehensive context for implementing Product CRUD with Clean Architecture in Go. It includes all necessary patterns, validation loops, specific implementation guidance, and complete context needed for successful one-pass implementation. The high confidence level is due to the detailed technical specifications, clear task breakdown, executable validation steps, and extensive research-backed best practices included.
# PRP - Store CRUD Implementation

**Purpose**  
Implement a complete Store CRUD API following Clean Architecture principles with PostgreSQL persistence, comprehensive validation, and testing. This PRP provides all necessary context for one-pass implementation success by leveraging existing Product CRUD patterns.

---

## Core Principles
- **Context is King**: Include all relevant docs, references, and examples needed for implementation.  
- **Validation Loops**: Provide automated build/tests/linting to ensure correctness at each stage.  
- **Information Dense**: Use project-specific patterns, naming, and conventions.  
- **Progressive Success**: Implement in minimal working form first, then enhance.  
- **Follow Architecture Rules**: Maintain strict dependency boundaries between layers.

---

## Goal
**Implement a complete Store CRUD API with PostgreSQL persistence following Clean Architecture**

Create endpoints for Create, Read (by ID and list), Update, and Delete operations for Store entities. The system must handle the specified database schema with proper validation, error handling, and follow Clean Architecture dependency rules, mirroring the existing Product CRUD implementation patterns.

---

## Why
- **Business Value**: Enables store management functionality as the foundation for multi-tenant product management system.  
- **Integration**: Serves as the missing piece that allows Products to reference Stores via store_id foreign key relationship.  
- **Problem Solved**: Provides structured, maintainable store data management with proper separation of concerns.

---

## What
**User-visible behavior:**  
- `POST /api/v1/stores` - Create new store with validation
- `GET /api/v1/stores/:id` - Retrieve single store by ID  
- `GET /api/v1/stores` - List all stores with optional pagination
- `PUT /api/v1/stores/:id` - Update existing store with validation
- `DELETE /api/v1/stores/:id` - Delete store by ID
- Proper HTTP status codes and JSON error responses
- Input validation with meaningful error messages

**Technical requirements:**  
- PostgreSQL table: `stores` with fields (id, name, description, address, phone)
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
- `url:` https://go.dev/doc/database/  
  `why:` Database access patterns, context handling, connection pooling
- `url:` https://pkg.go.dev/database/sql  
  `why:` QueryRowContext, ExecContext, transaction handling, error patterns
- `url:` https://gin-gonic.com/en/docs/examples/binding-and-validation/  
  `why:` Request binding, validation tags, error handling patterns  
- `url:` https://pkg.go.dev/github.com/stretchr/testify/mock  
  `why:` Mock creation, assertion patterns, table-driven tests
- `file:` `/internal/domain/product.go`  
  `why:` Domain entity pattern, validation methods, struct tags
- `file:` `/internal/usecase/product_usecase.go`  
  `why:` Business logic orchestration, logging, error wrapping patterns
- `file:` `/internal/repository/postgres/product_repository.go`  
  `why:` PostgreSQL implementation, parameterized queries, error handling
- `file:` `/internal/delivery/http/handlers/product_handler.go`  
  `why:` HTTP handler patterns, context timeouts, error response mapping
- `file:` `/internal/delivery/http/dto/product_dto.go`  
  `why:` Request/Response DTO patterns, validation tags, domain conversion
- `file:` `/migrations/001_create_products_table.up.sql`  
  `why:` PostgreSQL table creation, indexing patterns, timestamp handling

**Current Codebase Tree**  
```
/
├── cmd/
│   └── main.go                    # Application entry point with DI
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
```

**Desired Codebase Tree**  
```
/
├── internal/
│   ├── domain/
│   │   ├── store.go               # NEW: Store entity with business rules
│   │   ├── product.go             # EXISTING: Product entity 
│   │   └── errors.go              # MODIFY: Add store-specific errors
│   ├── usecase/
│   │   ├── interfaces.go          # MODIFY: Add StoreRepository interface
│   │   ├── store_usecase.go       # NEW: Store business logic orchestration
│   │   ├── store_usecase_test.go  # NEW: Unit tests with mocks
│   │   ├── product_usecase.go     # EXISTING: Product business logic
│   │   └── product_usecase_test.go # EXISTING: Product unit tests
│   ├── repository/
│   │   └── postgres/
│   │       ├── store_repository.go      # NEW: PostgreSQL implementation
│   │       ├── store_repository_test.go # NEW: Integration tests
│   │       ├── product_repository.go    # EXISTING: Product PostgreSQL impl
│   │       └── product_repository_test.go # EXISTING: Product integration tests
│   └── delivery/
│       └── http/
│           ├── dto/
│           │   ├── store_dto.go          # NEW: Store Request/Response DTOs
│           │   └── product_dto.go        # EXISTING: Product DTOs
│           ├── handlers/
│           │   ├── store_handler.go      # NEW: Store HTTP handlers
│           │   ├── store_handler_test.go # NEW: Store handler tests
│           │   ├── product_handler.go    # EXISTING: Product HTTP handlers
│           │   └── product_handler_test.go # EXISTING: Product handler tests
│           ├── middleware/               # EXISTING: No changes needed
│           └── router.go                 # MODIFY: Add store routes
├── migrations/
│   ├── 002_create_stores_table.up.sql   # NEW: Store table schema
│   ├── 002_create_stores_table.down.sql # NEW: Store table rollback
│   ├── 001_create_products_table.up.sql # EXISTING: Product schema
│   └── 001_create_products_table.down.sql # EXISTING: Product rollback
└── cmd/
    └── main.go                    # MODIFY: Add store dependency injection
```

**Known Gotchas**
- All repository interfaces must be defined in the use case layer (`internal/usecase/interfaces.go`), not repository layer
- Use parameterized queries ($1, $2) with lib/pq, never string concatenation
- Gin context timeout should be propagated to database operations using `context.WithTimeout`
- Use `database/sql.NullString` for nullable fields like description
- All HTTP handlers must validate input using struct tags before calling use cases
- Error responses should be consistent JSON format across all endpoints
- PostgreSQL `TEXT` type requires no length specification, `VARCHAR(n)` requires length
- Always use `RETURNING` clause in INSERT/UPDATE queries to get updated timestamps
- Mock interfaces must match exactly - use pointer receivers for repository methods
- Router integration requires updating both route registration and dependency injection in main.go

---

## Implementation Blueprint

### Database Schema & Migration
```sql
-- migrations/002_create_stores_table.up.sql
CREATE TABLE IF NOT EXISTS stores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    address TEXT,
    phone VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_stores_name ON stores(name);
CREATE INDEX idx_stores_created_at ON stores(created_at);

-- migrations/002_create_stores_table.down.sql  
DROP TABLE IF EXISTS stores;
```

### Data Models & Structure

**Domain Layer** (`internal/domain/store.go`):
- `Store` entity with business validation methods
- Validation rules: name required (1-100 chars), phone format validation, address/description max lengths
- Domain-specific errors in `errors.go` (ErrStoreNotFound, ErrInvalidStore, ErrDuplicateStore)

**Use Case Layer** (`internal/usecase/`):
- StoreRepository interface in `interfaces.go`  
- StoreUseCase struct with CRUD methods in `store_usecase.go`
- Business logic orchestration following product patterns
- Input validation and error logging

**Repository Layer** (`internal/repository/postgres/store_repository.go`):
- StoreRepository implementation
- Database connection handling with context propagation
- SQL query implementations with proper error handling
- Transaction support for complex operations

**Delivery Layer** (`internal/delivery/http/`):
- HTTP handlers for each CRUD operation in `store_handler.go`
- Request/Response DTOs with validation tags in `dto/store_dto.go`
- Route definitions in `router.go` under `/api/v1/stores`
- Error response formatting consistent with product handlers

### Task List

**Task 1: Domain Layer Implementation**  
CREATE `internal/domain/store.go`:
```go
type Store struct {
    ID          int64          `json:"id" db:"id"`
    Name        string         `json:"name" db:"name"`
    Description sql.NullString `json:"description" db:"description"`
    Address     sql.NullString `json:"address" db:"address"`
    Phone       sql.NullString `json:"phone" db:"phone"`
    CreatedAt   time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

func (s *Store) Validate() error
func (s *Store) IsValidPhoneFormat() bool
```

MODIFY `internal/domain/errors.go`:
- Add `ErrStoreNotFound`, `ErrInvalidStore`, `ErrDuplicateStore`

**Task 2: Use Case Layer Implementation**  
MODIFY `internal/usecase/interfaces.go`:
```go
type StoreRepository interface {
    Create(ctx context.Context, store *domain.Store) (*domain.Store, error)
    GetByID(ctx context.Context, id int64) (*domain.Store, error)
    GetAll(ctx context.Context, limit, offset int) ([]*domain.Store, error)
    Update(ctx context.Context, id int64, store *domain.Store) (*domain.Store, error)
    Delete(ctx context.Context, id int64) error
}

type StoreUseCaseInterface interface {
    CreateStore(ctx context.Context, store *domain.Store) (*domain.Store, error)
    GetStore(ctx context.Context, id int64) (*domain.Store, error)
    GetStores(ctx context.Context, limit, offset int) ([]*domain.Store, error)
    UpdateStore(ctx context.Context, id int64, store *domain.Store) (*domain.Store, error)
    DeleteStore(ctx context.Context, id int64) error
}
```

CREATE `internal/usecase/store_usecase.go`:
- StoreUseCase struct with repository dependency
- Implement all CRUD operations with business logic
- Follow logging and error handling patterns from `product_usecase.go`

**Task 3: Repository Implementation**  
CREATE `internal/repository/postgres/store_repository.go`:
- Implement StoreRepository interface
- Use parameterized queries with lib/pq following product patterns
- Handle PostgreSQL-specific errors with proper error mapping
- Include database transaction handling for complex operations

**Task 4: HTTP Delivery Layer**  
CREATE `internal/delivery/http/dto/store_dto.go`:
```go
type CreateStoreRequest struct {
    Name        string `json:"name" binding:"required,min=1,max=100"`
    Description string `json:"description" binding:"max=1000"`
    Address     string `json:"address" binding:"max=500"`
    Phone       string `json:"phone" binding:"max=100"`
}

type UpdateStoreRequest struct {
    Name        string `json:"name" binding:"required,min=1,max=100"`
    Description string `json:"description" binding:"max=1000"`
    Address     string `json:"address" binding:"max=500"`
    Phone       string `json:"phone" binding:"max=100"`
}

type StoreResponse struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Address     string `json:"address"`
    Phone       string `json:"phone"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type StoreListResponse struct {
    Stores []StoreResponse `json:"stores"`
    Total  int             `json:"total"`
    Limit  int             `json:"limit"`
    Offset int             `json:"offset"`
}
```

CREATE `internal/delivery/http/handlers/store_handler.go`:
```go
type StoreHandler struct {
    storeUseCase usecase.StoreUseCaseInterface
    logger       *logrus.Logger
}

func (h *StoreHandler) CreateStore(c *gin.Context)
func (h *StoreHandler) GetStore(c *gin.Context)
func (h *StoreHandler) GetStores(c *gin.Context)
func (h *StoreHandler) UpdateStore(c *gin.Context)
func (h *StoreHandler) DeleteStore(c *gin.Context)
```

**Task 5: Database Migration**  
CREATE `migrations/002_create_stores_table.up.sql` and `.down.sql`
- Include proper indexes for performance (name, created_at)
- Handle PostgreSQL-specific data types
- Ensure compatibility with existing product migrations

**Task 6: Router Integration**  
MODIFY `internal/delivery/http/router.go`:
```go
func SetupRouter(productHandler *handlers.ProductHandler, storeHandler *handlers.StoreHandler, logger *logrus.Logger) *gin.Engine {
    // ... existing setup ...
    
    api := r.Group("/api/v1")
    {
        // ... existing product routes ...
        
        stores := api.Group("/stores")
        {
            stores.POST("", storeHandler.CreateStore)
            stores.GET("/:id", storeHandler.GetStore)
            stores.GET("", storeHandler.GetStores)
            stores.PUT("/:id", storeHandler.UpdateStore)
            stores.DELETE("/:id", storeHandler.DeleteStore)
        }
    }
    return r
}
```

**Task 7: Application Wiring**  
MODIFY `cmd/main.go`:
- Add store repository, use case, and handler initialization
- Update router setup to include store handler
- Maintain dependency injection pattern from product implementation

**Task 8: Comprehensive Testing**  
CREATE `internal/usecase/store_usecase_test.go`:
- Mock repository tests following `product_usecase_test.go` patterns
- Table-driven tests with testify mocks
- Test all CRUD operations and error scenarios

CREATE `internal/delivery/http/handlers/store_handler_test.go`:
- HTTP handler tests following `product_handler_test.go` patterns
- Test all endpoints with proper mocking
- Validate request/response formats and error handling

CREATE `internal/repository/postgres/store_repository_test.go`:
- Integration tests following `product_repository_test.go` patterns
- Test database operations with real PostgreSQL connection
- Include edge cases and error scenarios

---

## Integration Points

**DATABASE**  
- PostgreSQL connection via lib/pq driver
- Migration management using golang-migrate
- Connection pooling via database/sql
- Transaction handling for atomic operations

**CONFIG**  
Environment variables (existing `.env`):
```env
# Existing configuration works for stores - no changes needed
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
- RESTful endpoint design: `/api/v1/stores`
- Consistent HTTP methods (POST, GET, PUT, DELETE)
- Proper HTTP status codes (200, 201, 400, 404, 500)

**DEPENDENCIES**  
Required Go modules (already available):
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
- `internal/usecase/store_usecase_test.go` - Mock repository tests
- `internal/delivery/http/handlers/store_handler_test.go` - HTTP handler tests  
- `internal/repository/postgres/store_repository_test.go` - Integration tests

**Level 3: Integration Tests**
```bash
# Database migration
migrate -path migrations -database "postgres://app_user:app_password@localhost:5432/product_db?sslmode=disable" up

# API testing with curl
curl -X POST http://localhost:8080/api/v1/stores \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Store", "description": "A test store", "address": "123 Main St", "phone": "+1234567890"}'

curl http://localhost:8080/api/v1/stores
```

---

## Testing Patterns

**Unit Tests with Mocks:**
```go
func TestStoreUseCase_CreateStore(t *testing.T) {
    tests := []struct {
        name    string
        store   *domain.Store
        mockFn  func(*mocks.StoreRepository)
        want    *domain.Store
        wantErr bool
    }{
        {
            name: "successful creation",
            store: &domain.Store{Name: "Test Store", Description: sql.NullString{String: "Test Description", Valid: true}},
            mockFn: func(m *mocks.StoreRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(&domain.Store{ID: 1}, nil)
            },
            want: &domain.Store{ID: 1},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &mocks.StoreRepository{}
            tt.mockFn(repo)
            
            uc := usecase.NewStoreUseCase(repo, logger)
            got, err := uc.CreateStore(context.Background(), tt.store)
            
            assert.Equal(t, tt.wantErr, err != nil)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**Integration Tests:**
```go  
func TestStoreRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewStoreRepository(db, logger)
    
    // Test Create
    store := &domain.Store{Name: "Integration Test Store", Description: sql.NullString{String: "Test Description", Valid: true}}
    created, err := repo.Create(context.Background(), store)
    
    require.NoError(t, err)
    assert.NotZero(t, created.ID)
    assert.Equal(t, store.Name, created.Name)
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
- [ ] Store creation, retrieval, update, and deletion work end-to-end
- [ ] Router integration includes all store endpoints
- [ ] Dependency injection in main.go updated for store components

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
❌ **Missing database indexes** - Index frequently queried fields (name, created_at)  
❌ **Ignoring context cancellation** - Propagate context through all layers  

---

## Performance Considerations
- Use database connection pooling (default with database/sql)
- Add indexes on name and created_at for frequently queried fields  
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

This PRP provides comprehensive context for implementing Store CRUD following the exact patterns established by the existing Product CRUD implementation. It includes all necessary patterns, validation loops, specific implementation guidance, and complete context needed for successful one-pass implementation. The high confidence level is due to the detailed technical specifications leveraging proven patterns, clear task breakdown, executable validation steps, and extensive research-backed best practices included.
# Go Clean Architecture — Service Template

A maintainable Go service structure based on **Clean Architecture** principles.  
This template uses a layered design for separation of concerns, testability, and adaptability.

---

## Folder Structure
```
/cmd/                # Application entry points (composition root)
/config/             # Configuration files and environment loading
/docs/               # API docs, ADRs, design documents
/internal/           # Application core (not exposed outside the module)
  /delivery/         # Transport layer (HTTP/gRPC handlers, DTOs, routers)
  /domain/           # Entities, value objects, domain services, business rules
  /repository/       # Repository interfaces & infrastructure implementations
  /usecase/          # Application services/use cases (business orchestration)
/migrations/         # Database migration scripts
/pkg                 # Imprement interface like gin, kafka-go
```

---

## Layer Responsibilities

**Domain**  
- Core business rules, invariants, and logic.  
- No external dependencies.  
- Entities, value objects, and domain-specific errors.

**Use Case**  
- Orchestrates business flows.  
- Calls repositories and services via interfaces.  
- Contains application-specific rules, not infrastructure.

**Repository**  
- Defines repository interfaces in `internal/repository`.  
- Implementations interact with DB or external services.  
- Keeps persistence logic separate from business rules.

**Delivery**  
- Handles incoming requests (HTTP, gRPC, messaging).  
- Maps DTOs to domain objects and vice versa.  
- Applies input validation and error mapping.

**Config**  
- Loads and validates environment variables and configurations.  
- Supports `.env` for development.

**Migrations**  
- SQL or tool-based migration files for schema changes.  
- Versioned and forward-only.

---

## Dependency Rules
- **Domain** is pure and does not depend on other layers.
- **Use Case** depends only on **Domain** and repository interfaces.
- **Delivery** depends on **Use Case** and DTOs.
- **Repository implementations** depend on DB clients or external APIs, not on **Delivery**.

---

## Example `.env.example`
```env
APP_NAME=service
APP_ENV=dev
HTTP_ADDR=0.0.0.0
HTTP_PORT=8080

DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=app_db
DB_SSLMODE=disable

LOG_LEVEL=info
```

---

## Request Flow
```
Client → Delivery (handler) → Use Case → Domain → Repository → Infrastructure → Response
```

---

## Testing Strategy
- **Unit tests**: domain entities/services, use cases with mock repositories.

---
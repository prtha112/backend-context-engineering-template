# Create PRP

## Feature file: $ARGUMENTS

Generate a complete PRP for general feature implementation with thorough research. Ensure context is passed to the AI agent to enable self-validation and iterative refinement. Read the feature file first to understand what needs to be created, how the examples provided help, and any other considerations.

The AI agent only gets the context you are appending to the PRP and training data. Assuma the AI agent has access to the codebase and the same knowledge cutoff as you, so its important that your research findings are included or referenced in the PRP. The Agent has Websearch capabilities, so pass urls to documentation and examples.

## Research Process

1. **Codebase Analysis**
    - Search for similar features/patterns
        - Controllers/Handlers (HTTP/gRPC) that match the feature domain
        - Reusable middlewares, request/response DTOs, validation patterns
        - Error handling conventions (error wrapping, sentinel errors, status mapping)
        - Config management (env, Viper, koanf, etc.)
        - DI/Composition patterns (wire, fx, dig, manual)
        - Persistence layer (sqlc/GORM/Ent/pgx), migrations (goose, migrate)
        - Caching/Queue patterns (redis, nats/kafka)
        - Logging/Tracing (zap, zerolog, slog), observability (otel)
        - Test layout and helpers (testify, httptest, golden files, mockery)
    - Identify files to reference in PRP
        - List absolute paths (e.g., internal/http/handlers/user.go, pkg/db/repo/order_repo.go)
    - Conventions to follow
        - Package layout (e.g., /cmd, /internal, /pkg)
        - Lint rules (golangci-lint config), CI steps, Make targets
        - Error model & response schema (problem+json? custom envelope?)
    - Test patterns
        - How unit/integration tests are structured
        - Mock strategy (interfaces + mockery? hand-written fakes?)

2. **External Research**
   - Library/Stdlib docs (URLs with specific sections)
        - Go: net/http, context, errors, database/sql / sync, time
        - Framework in use (detect from codebase): Gin/Chi/Echo/Fiber/grpc-go
        - ORM/Query: sqlc/GORM/Ent/pgx
        - Migrations: goose/golang-migrate
        - Config: Viper/koanf
        - Logging: zap/slog
        - DI: wire/fx/dig (or document manual wiring)
    - Implementation examples
        - Minimal handlers, middleware, transactional repo patterns
        - Test examples: httptest, table-driven tests, mocks
    - Best practices & pitfalls
        - Context timeouts/cancellation propagation
        - Don’t leak goroutines; defer rows.Close()
        - DB transaction boundaries; ReadCommitted vs app-level idempotency
        - JSON tag consistency; zero values; timezones
        - Race conditions & -race in tests

3. **User Clarification** (if needed)
   - Specific patterns to mirror and where to find them?
   - Integration requirements and where to find them?

## PRP Generation

Using PRPs/templates/prp_base.md as template:

### Critical Context to Include and pass to the AI agent as part of the PRP
- **Documentation**: URLs with specific sections
- **Code Examples**: Real snippets from codebase
- **Gotchas**: Library quirks, version issues
- **Patterns**: Existing approaches to follow

### Implementation Blueprint
- Start with pseudocode showing approach
- Reference real files for patterns
- Include error handling strategy
- list tasks to be completed to fullfill the PRP in the order they should be completed

### Validation Gates (Must be Executable) 
```bash
# Formatting
gofmt -l -w .
go mod tidy

# Static analysis
golangci-lint run
go vet ./...
staticcheck ./...

# Tests
go test ./... -v -race -count=1
```

*** CRITICAL AFTER YOU ARE DONE RESEARCHING AND EXPLORING THE CODEBASE BEFORE YOU START WRITING THE PRP ***

*** ULTRATHINK ABOUT THE PRP AND PLAN YOUR APPROACH THEN START WRITING THE PRP ***

## Output
Save as: `PRPs/{feature-name}.md`

## Quality Checklist
- [ ] All necessary context included
- [ ] Validation gates are executable by AI
- [ ] References existing patterns
- [ ] Clear implementation path
- [ ] Error handling documented

Score the PRP on a scale of 1-10 (confidence level to succeed in one-pass implementation using claude codes)

Remember: The goal is one-pass implementation success through comprehensive context.
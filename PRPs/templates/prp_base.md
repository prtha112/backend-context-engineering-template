# PRP Template – Go Clean Architecture (Context-Rich with Validation Loops)

**Purpose**  
Template optimized for Go projects using Clean Architecture to implement features with sufficient context and self-validation capabilities for delivering working, maintainable code through iterative refinement.

---

## Core Principles
- **Context is King**: Include all relevant docs, references, and examples needed for implementation.  
- **Validation Loops**: Provide automated build/tests/linting to ensure correctness at each stage.  
- **Information Dense**: Use project-specific patterns, naming, and conventions.  
- **Progressive Success**: Implement in minimal working form first, then enhance.  
- **Follow Architecture Rules**: Maintain strict dependency boundaries between layers.

---

## Goal
**[Clearly define the feature to be built, the desired end state, and functional outcomes.]**

---

## Why
- **Business Value**: [Describe why this feature matters and its user impact.]  
- **Integration**: [Explain how it fits with existing architecture/modules.]  
- **Problem Solved**: [Describe pain points and for whom they are solved.]

---

## What
- **User-visible behavior**: [Describe expected behavior from API/UI.]  
- **Technical requirements**: [List endpoints, data flows, constraints.]

---

## Success Criteria
- [List measurable acceptance criteria, e.g., “New endpoint responds <200ms for 95% of requests”]
- [Tests and linters pass]
- [All layers maintain dependency rules]

---

## All Needed Context
**Documentation & References**
- `url:` [Official library or API docs]  
  `why:` [Sections/methods you’ll use]
- `file:` [path/to/existing_handler.go]  
  `why:` [Pattern to follow and conventions to preserve]
- `doc:` [Internal guidelines URL or ADR]  
  `section:` [Relevant section about Clean Architecture usage]
  `critical:` [Pitfalls to avoid]

**Current Codebase Tree**  
[Paste output of `tree` for context]

**Desired Codebase Tree**  
[List new/modified files with purpose in each layer]

**Known Gotchas**
- Example: “All repositories must implement the interface in `internal/repository`.”
- Example: “All HTTP handlers must use `ctx` with request-scoped logging.”
- Example: “All DB changes require matching SQL in `/migrations`.”

---

## Implementation Blueprint
### Data Models & Structure
- **Domain Layer**: Add/modify entities and value objects.
- **Use Case Layer**: Define input/output ports and orchestrate logic.
- **Repository Layer**: Interface definitions and infrastructure implementations.
- **Delivery Layer**: HTTP handlers, request/response DTOs, validation.

### Task List
**Task 1:**  
MODIFY `internal/domain/<entity>.go`  
- Add new fields with validation methods.

**Task 2:**  
CREATE `internal/usecase/<feature>.go`  
- Follow pattern from `internal/usecase/<similar>.go`  
- Keep error handling consistent.

**Task 3:**  
MODIFY `internal/repository/<repo>.go`  
- Implement interface methods for new feature.  
- Maintain transaction and error-handling patterns.

**Task 4:**  
CREATE `internal/delivery/http/<feature>_handler.go`  
- Bind to router in `internal/delivery/http/router.go`.  
- Use existing request validation pattern.

**Task 5:**  
ADD migration in `/migrations` if schema changes are needed.

---

## Integration Points
**DATABASE**  
- Migrations & indexes as needed.

**CONFIG**  
- Add new config keys to `/config` with defaults from `.env`.

**ROUTES**  
- Add to HTTP router with appropriate path and method.

---

## Validation Loop

**Level 1: Syntax & Style**  
```bash
go fmt ./...
go vet ./...
golangci-lint run
```
Expected: No errors/warnings.

**Level 2: Unit Tests**  
- Create `internal/usecase/<feature>_test.go` and `internal/delivery/http/<feature>_handler_test.go`.  
- Follow table-driven test patterns from existing tests.

---

## Final Validation Checklist
- [ ] All tests pass: `go test ./...`
- [ ] No lint errors: `golangci-lint run`
- [ ] Manual API call works as expected
- [ ] Errors handled gracefully
- [ ] Logs are structured and contextual
- [ ] Documentation updated

---

## Anti-Patterns to Avoid
❌ Breaking dependency rules (outer layers importing inner business logic)  
❌ Skipping validation in DTOs or domain  
❌ Hardcoding values instead of using config  
❌ Catch-all `panic` recover without logging  
❌ Adding unused abstractions  

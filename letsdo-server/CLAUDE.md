## Project Overview

{{.ProjectDescription}}

**IMPORTANT**: This document contains general coding standards and patterns (HOW we build). For project requirements and specifications (WHAT we're building), refer to `SPECIFICATION.md`. For testing support, use `AGENTS.md` with the `/test` command.

## Quick Commands

```bash
make dev          # Run with hot reload (air)
make test         # Run all tests
make lint         # Run linter
make build        # Build binary
docker-compose up # Start dependencies (DB, Redis, etc.)
```

## Reference Docs

Before implementing features that touch database, APIs, or integrations, check the relevant doc in `docs/`:

| Task | Read First |
|------|------------|
| Database migrations, models, queries | `docs/schema/*.md` |
| API endpoint implementation | `docs/api/*.md` |
| Payment flows | `docs/integrations/payment.md` |
| SMS/Email notifications | `docs/integrations/notifications.md` |
| Auth/JWT implementation | `docs/integrations/auth.md` |
| Third-party services | `docs/integrations/*.md` |
| Architecture decisions | `docs/architecture/*.md` |
| Generated plans | `docs/plans/*.md` |

**Do NOT guess field types, constraints, or API shapes.** The docs have the exact specs from the design document.

The docs/schema/ folder contains one file per table with exact field specs, types, constraints, indexes, and business rules extracted from the design document.

If a doc doesn't exist for what you're implementing, ask before proceeding. If there are very sensible defaults for a table or api then you can use that while asking for confirmation.

## Project Structure

```
project-root/
├── cmd/
│   └── main.go              # Entrypoint only — wiring happens in config/
├── config/
│   ├── database.go          # DB connection + migrations
│   └── routes.go            # Route registration + dependency injection
├── internal/
│   ├── {domain}/            # One folder per domain (user, payment, order, etc.)
│   │   ├── model.go         # Entity + validation
│   │   ├── repository.go    # Data access interface + implementation
│   │   ├── service.go       # Business logic interface + implementation
│   │   └── handler.go       # HTTP handlers (defines its own Service interface)
│   ├── interfaces/          # Shared interfaces used across domains
│   ├── middleware/          # HTTP middleware
│   └── util/                # Generic helpers
├── migrations/              # SQL migrations (sequential numbering)
├── .gitignore/              # files to ignore, make sure we include binary build files and the like
├── .air.toml                # for air hotreloading when using make dev, make sure to configure to match where our main.go is, its NOT what is default init
└── Makefile
```

### Domain Package Pattern

Each domain is self-contained. Handler defines what it needs from Service. Service defines what it needs from Repository. This follows ISP — consumers own their interfaces.

```
internal/payment/
├── model.go         # Payment, Subscription entities
├── repository.go    # Repository interface + *repository implementation
├── service.go       # Service interface + *service implementation
├── handler.go       # Handler struct + Service interface it consumes
├── processor.go     # PaymentProcessor interface (for Stripe abstraction)
└── stripe.go        # Stripe implementation of PaymentProcessor
```

## Code Style (CRITICAL)

- Use `slog` for structured logging (NOT `log`)
- Error handling: wrap with `fmt.Errorf("context: %w", err)`
- Context propagation: always pass `ctx context.Context` as first param
- Naming: follow Go conventions (no get/set prefixes, no stuttering)

## Interface & Dependency Design

- **Define interfaces at point of use**, not implementation (consumer owns the interface)
- **Follow ISP**: interfaces expose only what the consumer needs
- **Follow DIP**: depend on interfaces, not concrete types
- **Follow IoC**: receiver controls dependencies (inject via constructor, never create internally)

Example — handler defines what it needs from service:

```go
// internal/payment/handler.go
type Service interface {
    CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error)
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
    // Only methods handler actually uses — NOT the full service
}
```

Example — service defines what it needs from repository:

```go
// internal/payment/service.go
type Repository interface {
    Create(ctx context.Context, payment *Payment) error
    GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
    // Only methods service actually uses
}
```

Example — cross-domain dependency via narrow interface:

```go
// internal/payment/service.go
type PaymentUserService interface {
    GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
    UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, customerID string) error
    // Only what payment service needs from user service
}
```

## Testing Requirements

- **TDD workflow**: Write failing tests first, then implement
- Test files live next to implementation: `service.go` → `service_test.go`
- Use table-driven tests for multiple cases
- Integration tests use `internal/testutil/suite.go` for setup

```bash
make test                    # All tests
make test-payment           # Single domain
go test ./internal/payment -run TestSpecificFunction -v
```

## Git Workflow

- Commit messages: `type: description` (feat, fix, test, refactor, chore, docs)
- Commit tests separately from implementation when doing TDD
- Never commit code that fails `make lint && make test`

## Architecture Rules

- All external services (Stripe, SendGrid, etc.) must be behind interfaces
- Database queries only in repository layer
- Business logic only in service layer
- Handlers do: parse request, call service, format response — nothing else
- No domain logic in handlers
- **Refer to SPECIFICATION.md for project requirements and business rules**

## Error Handling (CRITICAL - BASELINE REQUIREMENT)

### Repository Layer - ALWAYS Use errorutils.AnalyzeDBErr

**CRITICAL**: ALL database errors in repository methods MUST be wrapped with `errorutils.AnalyzeDBErr()`. This is NON-NEGOTIABLE. The codebase includes `internal/utils/errorutils` specifically for this purpose.

```go
// CORRECT - Always do this in repository methods:
func (r *repository) Create(ctx context.Context, entity *Entity) error {
    query := `INSERT INTO entities ...`
    err := r.db.ExecContext(ctx, query, ...)
    if err != nil {
        return errorutils.AnalyzeDBErr(err) // ✅ ALWAYS use this
    }
    return nil
}

// WRONG - Never do this in repository methods:
func (r *repository) Create(ctx context.Context, entity *Entity) error {
    query := `INSERT INTO entities ...`
    err := r.db.ExecContext(ctx, query, ...)
    if err != nil {
        return err // ❌ NEVER return raw DB errors
    }
    return nil
}
```

**What errorutils.AnalyzeDBErr handles automatically:**
- `sql.ErrNoRows` → `ErrNotFound`
- Duplicate key violations → `ErrDuplicateResource`
- Constraint violations → `ErrConstraintViolation`
- Other errors → passed through unchanged

### Handler Layer - Handle Known Errors

Handlers should check for specific errors from errorutils and return appropriate HTTP status codes:

```go
func (h *Handler) CreateEntity(c *gin.Context) {
    entity, err := h.service.CreateEntity(c.Request.Context(), req)
    if err != nil {
        switch {
        case errors.Is(err, errorutils.ErrDuplicateResource):
            c.JSON(409, gin.H{"error": "Entity already exists"})
        case errors.Is(err, errorutils.ErrInvalidInput):
            c.JSON(400, gin.H{"error": err.Error()})
        case errors.Is(err, errorutils.ErrConstraintViolation):
            c.JSON(400, gin.H{"error": "Invalid data provided"})
        case errors.Is(err, errorutils.ErrNotFound):
            c.JSON(404, gin.H{"error": "Entity not found"})
        case errors.Is(err, errorutils.ErrForbidden):
            c.JSON(403, gin.H{"error": "Access denied"})
        case errors.Is(err, errorutils.ErrUnauthorized):
            c.JSON(401, gin.H{"error": "Authentication required"})
        default:
            h.logger.Error("failed to create entity", slog.String("error", err.Error()))
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
        return
    }
    c.JSON(201, entity)
}
```

### Service Layer - Business Logic Errors

Services can add context or return business-specific errors:

```go
func (s *service) ProcessPayment(ctx context.Context, amount float64) error {
    if amount <= 0 {
        return errorutils.ErrInvalidInput // Use standard errors when applicable
    }

    payment, err := s.repo.CreatePayment(ctx, amount)
    if err != nil {
        // Repository already used AnalyzeDBErr, just add context if needed
        return fmt.Errorf("payment processing failed: %w", err)
    }
    return nil
}
```

## Database Transaction Handling

The codebase includes `internal/utils/dbutils` with transaction helpers for atomic operations:

### Using ExecTx for Transactions

Use `dbutils.ExecTx` when you need to perform multiple database operations atomically:

```go
// In service layer
func (s *service) TransferFunds(ctx context.Context, from, to uuid.UUID, amount float64) error {
    return dbutils.ExecTx(ctx, s.db, func(tx *sqlx.Tx) error {
        // All operations use the same transaction
        if err := s.repo.DebitAccountTx(ctx, tx, from, amount); err != nil {
            return err // Automatic rollback
        }

        if err := s.repo.CreditAccountTx(ctx, tx, to, amount); err != nil {
            return err // Automatic rollback
        }

        if err := s.repo.CreateTransferRecordTx(ctx, tx, from, to, amount); err != nil {
            return err // Automatic rollback
        }

        return nil // Automatic commit
    })
}

// In repository layer - provide both regular and Tx versions
func (r *repository) Create(ctx context.Context, entity *Entity) error {
    query := `INSERT INTO entities ...`
    _, err := r.db.ExecContext(ctx, query, ...)
    return errorutils.AnalyzeDBErr(err)
}

func (r *repository) CreateTx(ctx context.Context, tx *sqlx.Tx, entity *Entity) error {
    query := `INSERT INTO entities ...`
    _, err := tx.ExecContext(ctx, query, ...)
    return errorutils.AnalyzeDBErr(err)
}

## Environment Variables

```bash
# Required
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapp

# Optional
PORT=8000
JWT_SECRET=xxx
```

## Common Patterns

### Creating a New Domain

1. Create folder: `internal/{domain}/`
2. Create files: `model.go`, `repository.go`, `service.go`, `handler.go`
3. Define interfaces at consumer level (handler defines Service interface, etc.)
4. Wire up in `config/routes.go`
5. Add migration if needed

### Adding External Service Integration

1. Define interface in the domain that uses it: `internal/payment/processor.go`
2. Create implementation: `internal/payment/stripe.go`
3. Inject via constructor in `config/routes.go`
4. Write integration tests with real service in test mode

## Skills

This project includes Claude Code skills for a structured development workflow:

| Command | What it does |
|---------|-------------|
| `/write-a-prd` | Write a PRD through codebase exploration, submit as GitHub issue |
| `/prd-to-issues` | Break a PRD into vertical-slice task issues with dependency graph |
| `/tdd` | Implement a task using RED → GREEN → REFACTOR TDD cycle |

**Recommended workflow:**
1. `/write-a-prd` — Define the feature as a PRD issue
2. `/prd-to-issues #N` — Break it into independently-grabbable task issues
3. `/tdd #N` — Implement each task with test-driven development

## What NOT To Do

- Don't use `log` package — use `slog`
- Don't create dependencies inside functions — inject via constructor
- Don't put business logic in handlers
- Don't define interfaces in the implementing package
- Don't use `panic` for error handling
- Don't skip tests
- Don't modify files outside the scope of the current task

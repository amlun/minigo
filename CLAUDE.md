# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**minigo** is a Go web application template using Clean Architecture/DDD principles. It provides a foundation for building REST APIs with authentication, database management, and common utilities.

**Tech Stack:**
- **Framework:** Gin (HTTP router)
- **ORM:** Bun with PostgreSQL (pgdriver)
- **Auth:** JWT via golang-jwt/jwt
- **Config:** Viper + gotenv (.env support)
- **Logging:** logrus
- **ID Generation:** Snowflake (bwmarrin/snowflake)

## Commands

### Development
```bash
# Run the application (starts on port 8808 by default)
go run ./cmd/server

# Or use make
make run
```

### Building
```bash
# Build binary to build/ directory
make build

# Clean build artifacts
make clean
```

### Testing
```bash
# Run all tests
make test

# Run tests in a specific package
go test -v ./internal/domain/service

# Run tests with coverage
go test -cover ./...
```

### Dependencies
```bash
# Tidy dependencies
make tidy

# Update all dependencies
go get -u ./...

# Vendor dependencies
go mod vendor
```

### Database Migrations
Database migrations are stored in `migrations/` directory. Apply them manually using psql or your preferred migration tool:
```bash
psql -U postgres -d dbname -f migrations/001_init.sql
```

## Architecture

The project follows **Clean Architecture/Domain-Driven Design** with clear layer separation:

### Directory Structure
```
cmd/server/          - Application entry point
internal/
  domain/            - Domain layer (business logic, entities, interfaces)
    entity/          - Domain models (User, etc.)
    repository/      - Repository interfaces
    service/         - Domain services
    errors/          - Domain-specific errors
  application/       - Application layer (use cases, orchestration)
    service/         - Application services (AuthService, UserService, etc.)
  infrastructure/    - Infrastructure layer (external concerns)
    repository/      - Repository implementations (BunUserRepository)
    config/          - Configuration management (Viper)
    auth/            - JWT token generation/validation
    tx/              - Transaction management
    dbctx/           - Database context propagation
    id/              - Snowflake ID generator
    logging/         - Logger setup
    oss/             - OSS/file storage services
  interfaces/        - Interface layer (HTTP, DTOs, middleware)
    http/
      handlers/      - HTTP request handlers
      router.go      - Route definitions
    dto/             - Data Transfer Objects
    middleware/      - Auth, CORS, error handling, logging, rate limiting
    response/        - Standardized response helpers
pkg/                 - Shared utilities (hash, time, strings)
migrations/          - SQL migration files
```

### Key Architectural Patterns

**1. Transaction Management**
- Use `tx.Manager.InTx()` for transactional operations
- The transaction is injected into context via `dbctx.WithDB()`
- Repositories extract DB/Tx from context via `dbctx.FromCtx()`, enabling transparent transaction usage
- Example in service:
  ```go
  txManager.InTx(ctx, func(ctx context.Context) error {
      // All repository calls within this function use the same transaction
      return userRepo.Update(ctx, user)
  })
  ```

**2. Repository Pattern**
- Domain layer defines repository interfaces (e.g., `domain/repository/user_repository.go`)
- Infrastructure layer provides implementations (e.g., `infrastructure/repository/user_repository.go`)
- Repositories always accept `context.Context` as first parameter
- Use `dbctx.FromCtx(ctx, r.DB)` to get the correct DB/Tx handle

**3. Error Handling**
- Domain errors defined in `internal/domain/errors/` with typed error structure
- Application errors in `internal/application/service/errors.go`
- Infrastructure errors converted via `ConvertQueryError()`, `ConvertExecError()`
- Middleware `ErrorHandlerMiddleware()` catches panics and formats responses
- Use `response.HandleError()` or `middleware.AbortWithError()` for unified error responses
- Error types: SystemError, BusinessError, ValidationError, AuthError, NotFoundError

**4. Response Format**
- Unified JSON response structure via `internal/interfaces/response`
- Success: `response.Ok(c, data)` or `response.OkWithPage(c, list, total, page, pageSize)`
- Error: `response.HandleError(c, err)` or `response.ErrorFromAppError(c, appErr)`
- All responses include: `success`, `code`, `message`, and optionally `data`

**5. Pagination & Filtering**
- Use `dto.PaginationRequest` for standard pagination parameters
- Query filters available in `pkg/query` package
- Example: `NewQueryBuilder().Where(...).Like(...).Paginate(page, size).Apply(query)`
- Pagination response via `response.OkWithPage()`

**6. Request Tracking**
- Each request gets a unique ID via `RequestIDMiddleware()`
- Request ID is returned in `X-Request-ID` header
- All logs include request_id for tracing
- Access via `middleware.GetRequestID(c)`

**7. Authentication Flow**
- JWT tokens generated in `infrastructure/auth/jwt.go`
- Middleware `AuthMiddleware()` validates tokens and injects user info into context
- `RoleAuthMiddleware()` enforces role-based access control
- User passwords hashed with bcrypt in entity hooks (`BeforeInsert`, `BeforeUpdate`)

**8. Configuration**
- Config loaded from `.env` file (via gotenv) and environment variables
- Viper provides unified access with `config.GetPort()`, `config.GetDBDsn()`, etc.
- Supports multiple config file formats (YAML, JSON, TOML) - searched in `./`, `./cmd/server`, `./backend`

**9. Router Setup**
- `BuildRouter(db *bun.DB)` in `internal/interfaces/http/router.go` wires up dependencies
- Dependencies injected top-down: repositories → services → handlers
- Middleware order: CORS → RequestID → ErrorHandler → Recovery → RequestLogger

**10. ID Generation**
- Snowflake IDs used for primary keys (64-bit integers)
- Initialized once at startup via `id.Init()` in main.go
- Entity IDs should be `int64` type with `bun:"id,pk,autoincrement"` tag (though Snowflake generates them, not DB)

## Configuration

Environment variables (`.env` file or system env):
```
ENV=dev|prod              # Environment mode (affects logging, Gin mode)
PORT=8808                 # Server port
DB_DSN=postgres://...     # PostgreSQL connection string
LOG_LEVEL=debug|info|warn|error
JWT_SECRET=...            # Secret key for JWT signing
JWT_EXPIRE_DURATION=24h   # Token expiration (duration format)
```

OSS configuration (if using object storage):
```
OSS_ENDPOINT=...
OSS_ACCESS_KEY_ID=...
OSS_ACCESS_KEY_SECRET=...
OSS_BUCKET_NAME=...
OSS_REGION=...
OSS_TOKEN_EXPIRE_SECONDS=3600
```

## Testing Strategy

- Unit tests live alongside code (e.g., `example_service_test.go`)
- Use `make test` or `go test ./...` to run all tests
- For repository tests, consider using testcontainers or an in-memory database

## Common Development Patterns

**Adding a New Entity:**
1. Define entity in `internal/domain/entity/` with Bun tags
2. Create repository interface in `internal/domain/repository/`
3. Implement repository in `internal/infrastructure/repository/`
4. Create migration SQL in `migrations/`
5. Add application service in `internal/application/service/`
6. Create DTOs in `internal/interfaces/dto/`
7. Add handler in `internal/interfaces/http/handlers/`
8. Register routes in `internal/interfaces/http/router.go`

**Adding a New Endpoint:**
1. Define DTOs for request/response in `internal/interfaces/dto/`
2. Create handler method in appropriate handler file
3. Use helper functions:
   - `middleware.ValidateAndBindJSON(c, &req)` for JSON body
   - `middleware.ValidateIDParam(c, "id")` for path parameters
   - `middleware.MustNotError(c, err)` to abort on error
4. Return responses:
   - Success: `response.Ok(c, data)`
   - Paginated: `response.OkWithPage(c, list, total, page, size)`
   - Error: `response.HandleError(c, err)`
5. Register route in `BuildRouter()` with necessary middleware
6. For authenticated routes, use `AuthMiddleware()` and optionally `RoleAuthMiddleware()`

**Implementing Business Logic:**
- Simple CRUD → Application service
- Complex domain rules → Domain service (in `internal/domain/service/`)
- Cross-entity operations → Application service with transaction management

**Pagination Query Example:**
```go
// In handler
var req dto.PageSortRequest
if !middleware.ValidateAndBindQuery(c, &req) {
    return
}

// In service/repository
list, total, err := repo.List(ctx, req.GetPage(), req.GetPageSize())
if err != nil {
    return nil, 0, err
}

// In handler response
response.OkWithPage(c, list, total, req.GetPage(), req.GetPageSize())
```

**Query Filtering Example:**
```go
import "minigo/pkg/query"

qb := query.NewQueryBuilder().
    Like("name", searchKeyword).
    DateRange("created_at", startTime, endTime).
    Order("created_at", true).
    Paginate(page, pageSize)

query := db.NewSelect().Model(&models)
query = qb.Apply(query)
err := query.Scan(ctx)
```

**Error Handling Example:**
```go
// In service layer
if user == nil {
    return apperrors.ErrUserNotFound
}

// In handler
if !middleware.MustNotError(c, err) {
    return
}

// Or manually
if err != nil {
    response.HandleError(c, err)
    return
}
```

**Validation Example:**
```go
import "minigo/pkg/validator"

// Use built-in validators
if !validator.IsPhone(phone) {
    return apperrors.NewValidationError("VAL_PHONE", "Invalid phone number")
}

if !validator.LengthBetween(username, 3, 20) {
    return apperrors.NewValidationError("VAL_LENGTH", "Username must be 3-20 characters")
}

// In DTOs, use struct tags
type CreateUserRequest struct {
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required,min=2,max=50"`
}
```

## Notes

- **Vendor directory:** Dependencies are vendored. Run `go mod vendor` after updating dependencies.
- **Deployment:** `deploy.sh` is for deployment automation (references cloud platform setup).
- **Default port:** 8808
- **Health check endpoint:** `GET /api/health`

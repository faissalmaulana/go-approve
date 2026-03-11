# AGENTS.md

## Project Overview

This is a Go (Golang) REST API project using:
- **Echo** - HTTP web framework
- **GORM** - ORM for MySQL database
- **FX** - Dependency injection framework
- **Zap** - Structured logging
- **Testify** - Testing toolkit with assertions

## Build Commands

### Run the application
```bash
go run ./cmd
```

### Build binary
```bash
go build -o ./tmp/main ./cmd
```

### Run with hot reload (air)
```bash
air
```
This uses `.air.toml` configuration. The binary is built to `./tmp/main`.

### Run tests
```bash
go test ./...
```

### Run a single test
```bash
go test -v ./internal/service/jwtfx/... -run TestGenerateAccessToken
```

### Test with coverage
```bash
go test -cover ./...
```

## Code Style Guidelines

### Naming Conventions

- **Files**: Use lowercase with underscores (snake_case), e.g., `user_repository.go`, `jwt_test.go`
- **Packages**: Use short, lowercase names, e.g., `service`, `repository`, `handlers`
- **Interfaces**: Use descriptive names ending with `er` or meaningful nouns, e.g., `UserStorage`, `Reader`
- **Variables**: Use camelCase, e.g., `userPayload`, `accessToken`
- **Constants**: Use PascalCase for exported, camelCase for unexported, e.g., `QueryTimeout`
- **Functions**: Use PascalCase for exported, camelCase for unexported
- **Structs**: Use PascalCase for both exported and unexported

### Imports

Group imports in the following order (blank line between groups):
1. Standard library packages
2. Third-party packages
3. Internal packages

```go
import (
    "context"
    "net/http"
    "time"

    "github.com/faissalmaulana/go-approve/cmd/handlers"
    "github.com/labstack/echo/v5"
    "go.uber.org/fx"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "github.com/faissalmaulana/go-approve/internal/db"
    "github.com/faissalmaulana/go-approve/internal/repository/user"
    "github.com/faissalmaulana/go-approve/internal/service/auth"
)
```

### Project Structure

```
cmd/
    main.go          # Application entry point
    mux.go            # HTTP router setup
    server.go         # HTTP server configuration
    handlers/         # HTTP request handlers
    dto/              # Data Transfer Objects

internal/
    constant/         # Constants
    model/            # Database models
    repository/       # Data access layer
    service/          # Business logic
    db/               # Database connection
```

### Error Handling

- Define errors as package-level variables using `errors.New()`
- Use sentinel errors in the `internal/service/errors.go` file
- Return errors explicitly, handle them in handlers with appropriate HTTP status codes

```go
// Define errors
var (
    ErrDuplicatedUser = errors.New("user already exists")
    ErrSubIsEmpty     = errors.New("jwt subject is required")
)

// Handle errors in handlers
if err != nil {
    switch err {
    case service.ErrDuplicatedUser:
        return echo.NewHTTPError(http.StatusConflict, err.Error())
    default:
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
}
```

### Dependency Injection

Use FX for dependency injection. Provide constructors that return interfaces:

```go
fx.Provide(
    NewEchoMux,
    db.New,
    fx.Annotate(user.New, fx.As(new(user.UserStorage))),
    handlers.NewHealthHandler,
)
```

### Testing

- Test files use `_test.go` suffix
- Use `testify/assert` for assertions
- Use subtests with `t.Run()` for better organization

```go
func TestGenerateAccessToken(t *testing.T) {
    j := JwtFx{...}

    t.Run("success generate access token", func(t *testing.T) {
        result, err := j.GenerateAccessToken("test")
        assert.NoError(t, err)
        assert.NotEmpty(t, result)
    })

    t.Run("error generate access token because subject is empty", func(t *testing.T) {
        result, err := j.GenerateAccessToken("")
        assert.ErrorIs(t, err, service.ErrSubIsEmpty)
        assert.Empty(t, result)
    })
}
```

### Database Models

- Use GORM tags for column definitions
- Implement `BeforeCreate` hook for UUID generation
- Use `gorm.DeletedAt` for soft deletes

```go
type User struct {
    ID        string         `gorm:"primaryKey"`
    Name      string         `gorm:"type:varchar(100);not null"`
    Handler   string         `gorm:"type:varchar(50);unique;not null"`
    Email     string         `gorm:"type:varchar(255);unique;not null"`
    Password  string         `gorm:"type:varchar(60);not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### Logging

Use Zap for structured logging:

```go
zap.NewProduction()
```

### Context Usage

Always pass context and use timeouts for database operations:

```go
ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
defer cancel()
```

### Configuration

- Use `.env` file for environment variables (loaded via `joho/godotenv`)
- Create config structs in respective service packages

### Best Practices

1. **Keep handlers thin** - delegate business logic to services
2. **Use interfaces** - define interfaces in package that uses them
3. **Return early** - avoid nested conditionals when possible
4. **Close resources** - always defer cleanup (e.g., context cancel)
5. **Use meaningful names** - avoid single-letter variables except in loops
6. **Comment exported functions** - add doc comments for public APIs
7. **Group related code** - keep related functions together in files

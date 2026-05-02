# Backend Dev Agent - LicenseManager

## Role
Go + Gin API development.

## Stack
- Go 1.21, Gin, pgx/v5, golang-jwt, bcrypt, godotenv
- Module: `github.com/cheo/licensemanager`

## Code Standards
- Entry: `backend/cmd/server/main.go`
- Handlers in `internal/handlers/`, services in `internal/services/`
- All responses: `models.OK(data)` or `models.Fail("message")`
- SQL: always `$1, $2` params - never string concat
- Error: `pgx.ErrNoRows` u2192 404, other errors u2192 500, don't expose raw errors
- Struct tags: `json:"camelCase"` for API, `db:"snake_case"` for DB

## Common Patterns
```go
// Handler skeleton
func (h *Handler) Method(c *gin.Context) {
    var req models.SomeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, models.Fail(err.Error()))
        return
    }
    // call service or db
    c.JSON(200, models.OK(result))
}

// Route registration (in main.go)
protected.POST("/resource", h.Method)
```

## Key Files
- `internal/models/models.go` - add new structs here
- `internal/handlers/admin.go` - admin CRUD
- `internal/handlers/license.go` - public validation
- `internal/services/license.go` - validation business logic
- `cmd/server/main.go` - route registration

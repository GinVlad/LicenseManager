# Backend Rules - LicenseManager

## Stack
- Language: Go 1.21+
- Framework: Gin
- Database: PostgreSQL via pgx/v5
- Auth: JWT (admin) + API Key header (client apps)

## Project Layout
```
backend/
  |-- cmd/server/main.go       # Entry: loads config, DB, runs Gin
  |-- internal/
      |-- config/config.go      # Env vars via godotenv
      |-- database/
          |-- postgres.go          # pgxpool.Pool connection
          |-- migrations/          # .up.sql / .down.sql
          |-- handlers/
              |-- admin.go             # Admin CRUD (apps, licenses, HWIDs)
              |-- license.go           # Public validate + remove HWID
          |-- middleware/
              |-- auth.go              # JWTAuth(), AppAPIKey()
              |-- ratelimit.go         # RateLimit(n, window)
              |-- security.go          # SecureHeaders(), CORS(origins)
          |-- models/models.go      # Structs, DTOs, OK()/Fail() helpers
          |-- services/license.go   # Validate(), GenerateLicenseKey()
          |-- go.mod
```

## Conventions
- Handlers call services, services call db directly (no repo layer - app is simple)
- All handlers return `models.APIResponse{ success, data, error }`
- Use `c.ShouldBindJSON` for request parsing, return 400 on error
- Use `pgx.ErrNoRows` to detect not-found vs real errors
- Admin routes: POST `/admin/login` (public) + JWT group
- License routes: protected by `AppAPIKey` middleware (X-App-Key header)

## Error Handling
```go
if err == pgx.ErrNoRows {
    c.JSON(404, models.Fail("not found"))
    return
}
if err != nil {
    c.JSON(500, models.Fail("server error"))
    return
}
```

## Adding a New Endpoint
1. Add struct to `models/models.go` if new shape needed
2. Add handler method to `handlers/admin.go` or `handlers/license.go`
3. Register route in `cmd/server/main.go`
4. Use `/generate-api-endpoint` skill for scaffolding

## Adding a New Migration
1. Create `internal/database/migrations/00N_description.up.sql`
2. Create matching `.down.sql`
3. Migrations run automatically on server start via `database.RunMigrations()`
4. Use `/generate-migration` skill

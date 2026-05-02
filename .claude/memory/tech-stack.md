---
name: Tech Stack
description: Exact versions, module paths, and code patterns for LicenseManager
type: project
---

## Backend
- Go 1.25, module: `github.com/cheo/licensemanager`
- Gin v1.10, pgx/v5, golang-jwt/jwt v5, bcrypt, godotenv v1.5
- Entry: `backend/cmd/server/main.go`
- Migrations: auto-run from `backend/internal/database/migrations/*.up.sql`

## Frontend
- React 18, Vite 5, Tailwind v3, React Router v6, Lucide React
- Proxy: Vite proxies `/api` u2192 `http://localhost:8080`
- Token: stored in `localStorage` as `lm_token`
- API client: `frontend/src/lib/api.js` - all calls throw Error on failure

## SDK
- `sdk/go/license/client.go` - import path: `github.com/cheo/licensemanager/sdk/go/license`
- HWID: SHA256(MAC + hostname + $USER)[:16] as hex = 32 chars
- Timeout: 15s per request

## Key File Paths
```
backend/cmd/server/main.go              # routes wired here
backend/internal/handlers/admin.go      # admin CRUD
backend/internal/handlers/license.go    # public validate
backend/internal/services/license.go    # validation business logic
backend/internal/models/models.go       # all structs + DTOs
backend/internal/middleware/auth.go     # JWTAuth(), AppAPIKey()
frontend/src/lib/api.js                 # all fetch calls
frontend/src/App.jsx                    # routes
frontend/src/components/Layout.jsx      # sidebar nav
sdk/go/license/client.go                # SDK for Wails apps
```

## How to apply
When modifying any of these files, read the file first - do not guess the current structure.
Module paths and import aliases match what's in go.mod exactly.

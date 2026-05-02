# LicenseManager

Self-hosted license management SaaS. One server manages licenses for multiple apps.
Admin dashboard to issue/revoke/extend licenses. Go SDK for client apps to validate.

## Quick Reference
- **Stack:** Go + Gin + PostgreSQL | React + Tailwind
- **Auth:** JWT (admin) + API Key (client apps)
- **Purpose:** Issue license keys, bind HWIDs, enforce expiry, multi-app support

## Project Layout
```
backend/          # Go + Gin API server
  cmd/server/     # Entry point
  internal/
    config/       # Env config loader
    database/     # PG pool + migrations
    handlers/     # HTTP handlers (admin, license)
    middleware/   # JWT, API key, rate limit, CORS, security headers
    models/       # Structs + DTOs
    services/     # Business logic (license validation, HWID)
frontend/         # React admin dashboard
  src/
    pages/admin/  # Dashboard, Apps, Licenses
    pages/auth/   # Login
    components/   # Layout, shared UI
    lib/api.js    # API client
sdk/go/license/   # Drop-in Go client for Wails/CLI apps
```

---

## Agentic Workflow

Use `/agentic [feature]` to start coordinated implementation:

```
BRAINSTORM → PLAN → COOK → VERIFY → FIX → DONE
```

| Command | Action |
|---------|--------|
| `/agentic [feature]` | Full flow with gates |
| `/agentic resume` | Continue from session.md |
| `/agentic brainstorm [feature]` | Single stage |

---

## Rules System

| Task | Rules to Load |
|------|---------------|
| Backend (Go/Gin) | `rules/backend.md` |
| Frontend (React) | `rules/frontend.md` |
| Database (PostgreSQL) | `rules/database.md` |
| Security | `rules/security.md` |
| Deployment | `rules/deployment.md` |
| Testing | `rules/testing.md` |

**Always check `rules/session.md` first.**

---

## Commands
```bash
# Dev
cp .env.example .env
docker compose up -d          # Start postgres + backend + frontend
open http://localhost:5173    # Admin dashboard

# Backend only
cd backend && go run ./cmd/server

# Frontend only
cd frontend && npm install && npm run dev
```

## Key Concepts
- **App** — a product you sell (e.g. "eBay Creator"). Has a slug + API key.
- **License** — a key bound to one app. Has plan, max threads, expiry, active flag.
- **HWID** — machine fingerprint. Each license can bind N HWIDs (set per app).
- **Validation flow** — app sends `{ appSlug, key, hwid }` → server checks in order:
  1. App API key valid
  2. License key exists for this app
  3. License is active
  4. Not expired
  5. HWID registered → update last_seen
  6. HWID slot available → register new device
  7. Return `{ valid, expiresAt, maxThreads }`

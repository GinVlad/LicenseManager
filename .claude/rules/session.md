# Session State

## Last Updated
2026-04-26

## Current Phase
**Phase 2: Stack running — testing client flow**

---

## Active Workflow
**Feature:** None
**Stage:** -
**Plan File:** -

### Workflow Progress
```
[ ] BRAINSTORM
[ ] PLAN
[ ] COOK
[ ] VERIFY
[ ] FIX (if needed)
[ ] DONE
```

---

## What's Done
- [x] PostgreSQL schema (admins, users, applications, licenses, license_hwids)
- [x] Go + Gin backend (config, database, middleware, handlers, services)
- [x] License validation with multi-HWID support
- [x] Admin handlers: apps CRUD, license issue/update/delete, HWID management
- [x] JWT auth middleware (admin) + API key middleware (client apps)
- [x] Rate limiting (120 req/min per IP), CORS, secure headers
- [x] React admin dashboard: Login, Dashboard, Apps, Licenses pages
- [x] Go SDK client (sdk/go/license/client.go)
- [x] docker-compose, Dockerfile, .env.example, README.md
- [x] Full stack running via docker compose (postgres:5434, backend:8080, frontend:5173)
- [x] Fixed Dockerfile Go version (1.21 -> 1.25 to match go.mod)
- [x] Fixed Vite proxy: BACKEND_URL=http://backend:8080 inside Docker network
- [x] Admin login confirmed working
- [x] Per-license HWID limit (`max_hwid` on licenses table, migration 002)
- [x] Migration runner upgraded: sequential, tracks applied versions in schema_migrations
- [x] Frontend: Max Machines field in issue form + Machines column in license table
- [x] Validated: Basic(maxHwid=1) blocks 2nd machine, Pro(maxHwid=3) allows 3 then blocks 4th

## What's Next
- [ ] Connect eBay Creator Wails app to this server (replace placeholder ValidateLicense)
- [ ] Add Stripe payment integration (optional)
- [ ] Add user self-service portal (optional)

---

## Blockers
None

## Decisions Made
| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | Gin | Middleware ecosystem (rate limit, JWT, CORS) |
| Database | PostgreSQL | Multi-app scalability, UUID PKs |
| Auth | JWT (admin) + API Key (apps) | Simple, stateless |
| HWID per license | Configurable per app | Flexible set 1 for strict, N for teams |

## Notes for Next Session
- Backend runs on :8080, frontend on :5173 (proxies /api to backend)
- Admin seeded from ADMIN_EMAIL + ADMIN_PASSWORD env vars on first boot
- License keys generated as `{APP_PREFIX}-{XXXX}-{XXXX}-{XXXX}`

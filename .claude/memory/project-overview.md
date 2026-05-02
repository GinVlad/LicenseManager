---
name: Project Overview
description: What LicenseManager is, its architecture, key design decisions, and how to extend it
type: project
---

LicenseManager is a self-hosted SaaS for managing software licenses across multiple desktop apps.
Built by cheo to license eBay Creator and future apps from one server.

**Why:** Needed license key enforcement with HWID binding, expiry, and revocation for Wails desktop apps.

**How to apply:** When extending the project, every decision should support multi-app (N apps, not just eBay).

## Architecture
```
Desktop App (Wails/Go)
  - sdk/go/license/client.go -> POST /api/v1/license/validate (X-App-Key)
                                   |
                          LicenseManager API (Gin :8080)
                                   |
                          PostgreSQL (applications, licenses, license_hwids)

Admin Browser -> React :5173 -> /api/v1/admin/* (JWT Bearer)
```

## Key Decisions
- **Gin over Chi**: middleware ecosystem (rate limit, JWT, CORS, security headers)
- **PostgreSQL over SQLite**: multi-app at scale, UUID PKs, concurrent writes
- **No ORM**: pgx direct queries - simple enough, full control
- **Multi-HWID per license**: `license_hwids` table, `max_hwid_per_license` set per app
- **API key per app**: `applications.api_key` - auto-generated, isolates apps from each other
- **Admin JWT 24h**: short enough to be secure, long enough to not annoy

## Current State (2026-04-26)
Phase 1 complete: full backend + admin frontend + Go SDK built.
Not yet connected to a live eBay Creator instance (placeholder ValidateLicense still in use).

## Planned Extensions
- User self-service portal (view my licenses, remove my devices)
- Stripe payments (checkout -> auto-issue license)
- Email notifications (expiry warning, new device registered)

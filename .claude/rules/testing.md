# Testing Rules - LicenseManager

## Backend (Go)

### Unit Tests
- Location: alongside source, `_test.go` suffix
- Focus: `services/license.go` - Validate logic is the core, test all branches
- Mock: use an in-memory map or test PostgreSQL (prefer real DB)

### Integration Tests
```bash
# Spin up test DB
docker run -d -p 5433:5432 -e POSTGRES_PASSWORD=test postgres:16-alpine

DATABASE_URL=postgres://postgres:test@localhost:5433/licensemanager?sslmode=disable \
  go test ./...
```

### Key Scenarios to Test
- Validate: valid key + new HWID -> registers device, returns valid
- Validate: valid key + known HWID -> updates last_seen, returns valid
- Validate: expired key -> returns { valid: false, error: "license expired" }
- Validate: revoked key -> returns { valid: false, error: "license revoked" }
- Validate: max HWIDs reached -> returns { valid: false, error: "max devices", devices: [...] }
- Validate: unknown key -> returns { valid: false, error: "invalid license key" }
- Admin login: wrong password -> 401
- Admin routes without JWT -> 401
- License routes without X-App-Key -> 401

## Frontend (React)

### Manual Test Checklist
- [ ] Login with wrong password shows error
- [ ] Login with correct password redirects to dashboard
- [ ] Dashboard shows correct counts
- [ ] Create app u2192 appears in list with API key
- [ ] Issue license u2192 appears in list with ACTIVE status
- [ ] Revoke license u2192 status changes to REVOKED
- [ ] Expand HWID row u2192 shows devices
- [ ] Delete HWID u2192 device removed from list
- [ ] Logout clears token and redirects to login

## SDK Test (Go)
```go
client := license.New(license.Config{
    ServerURL: "http://localhost:8080",
    AppSlug:   "test-app",
    AppAPIKey: "<key from admin>",
})
info, err := client.Validate("TEST-XXXX-XXXX-XXXX")
// err should be nil, info.Valid == true
```

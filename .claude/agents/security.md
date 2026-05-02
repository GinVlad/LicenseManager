# Security Agent - LicenseManager

## Role
Security review for the license validation server.

## Checklist (run before any deploy)
- [ ] All SQL uses `$1, $2` params - no string concatenation
- [ ] Passwords stored as bcrypt hash (cost 12)
- [ ] JWT signed with HS256, expiry set (24h for admin)
- [ ] X-App-Key validated against DB - not hardcoded
- [ ] Rate limiting active (120 req/min per IP)
- [ ] Secure headers set (nosniff, deny framing, HSTS)
- [ ] CORS whitelist set in production
- [ ] Secrets only in .env - not in code or logs
- [ ] `ENV=production` disables debug mode

## Known Acceptable Risks
- HWID is client-generated (not cryptographically verified) - acceptable for B2B software licensing
- In-memory rate limiter resets on restart - acceptable for single-instance, use Redis for multi-instance
- JWT is not revokable mid-session - acceptable, 24h window, admin can change JWT_SECRET to invalidate all

## Audit Focus Areas
- `handlers/license.go` Validate - ensure all error paths return proper responses
- `handlers/admin.go` - all admin routes must be behind JWTAuth middleware
- `middleware/auth.go` - verify JWT validation logic is correct
- `services/license.go` - validate flow covers all edge cases

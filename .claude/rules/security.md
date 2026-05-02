# Security Rules - LicenseManager

## Auth Model
```
Admin panel:   POST /admin/login -> JWT (24h)
               Protected routes: Authorization: Bearer <token>

Client apps:   X-App-Key: <api_key_from_applications_table>
               Validated by AppAPIKey() middleware per request
```

## Password Hashing
- bcrypt cost 12 (see `golang.org/x/crypto/bcrypt`)
- Never store plaintext. Never log passwords.

## JWT
- HS256, signed with JWT_SECRET env var
- Claims: `sub` (admin ID), `exp` (24h)
- Validated in `middleware/auth.go` JWTAuth()

## Rate Limiting
- 120 req/min per IP (in-memory bucket, resets each minute)
- Applies globally to all routes
- Returns 429 on breach

## Secure Headers (set by SecureHeaders middleware)
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

## CORS
- Configured in `middleware/security.go` CORS(origins)
- Whitelist specific origins in production via CORS_ORIGINS env
- Allows: Content-Type, Authorization, X-App-Key

## Input Validation
- Use `c.ShouldBindJSON(&req)` with `binding:"required"` tags
- Never interpolate user input into SQL - always use `$1, $2` params
- HWID and API keys are opaque strings - validate length if needed

## Secrets in .env (never commit)
```
DATABASE_URL   # includes password
JWT_SECRET     # at least 32 random bytes
ADMIN_PASSWORD # seeded on first boot
```

## Threat Model for License Validation
- **Key brute-force**: 12-byte random key = 96 bits entropy, rate limit prevents enumeration
- **HWID spoofing**: HWID is client-generated, not cryptographically verified - acceptable for software licensing
- **Replay attacks**: Each validation updates `last_validated_at` - not a secret, no replay risk
- **API key leak**: Rotate via admin panel (delete app + recreate). Keys are per-app, blast radius is one app.

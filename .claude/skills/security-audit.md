# Skill: Security Audit

## Usage
`/security-audit`

## What It Does
Runs through the LicenseManager security checklist and reports findings.

## Checklist

### SQL Injection
```bash
grep -rn 'Sprintf.*SELECT\|Sprintf.*INSERT\|Sprintf.*UPDATE\|Sprintf.*DELETE' backend/
# Should return 0 results
```

### Hardcoded Secrets
```bash
grep -rn 'password\|secret\|apikey\|token' backend/ --include='*.go' | grep -v '_test.go' | grep -v '.env'
# Review each result
```

### Route Middleware
```bash
cat backend/cmd/server/main.go
# Verify: all /admin routes (except /login) are inside the JWT group
# Verify: all /license routes are inside the AppAPIKey group
```

### CORS Origins
```bash
grep CORS_ORIGINS .env
# Should not be empty in production
```

### Rate Limit
```bash
grep RateLimit backend/cmd/server/main.go
# Should be applied globally
```

## Output
Report findings as: PASS / WARN / FAIL with remediation steps.

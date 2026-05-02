# DevOps Agent - LicenseManager

## Role
Docker, deployment, and infrastructure.

## Local Dev Stack
```yaml
docker compose up -d
  postgres:16-alpine  # :5432
  backend             # :8080 (Go + Gin)
  frontend            # :5173 (Vite dev server)
```

## Production Docker
```bash
# Multi-stage build - final image is alpine ~20MB
docker build -t licensemanager .
docker run -p 8080:8080 --env-file .env licensemanager
```

## Required Env Vars
```
PORT, ENV, DATABASE_URL, JWT_SECRET, ADMIN_EMAIL, ADMIN_PASSWORD
```

## Deployment Targets
- **Railway**: `railway init && railway add --database postgresql && railway up`
- **Fly.io**: `fly launch && fly postgres create && fly deploy`
- **VPS**: docker + nginx reverse proxy

## Health Check
```bash
curl http://localhost:8080/health
# { "status": "ok", "time": "..." }
```

## Frontend Deploy
```bash
cd frontend && npm run build
# Upload dist/ to Vercel, Netlify, or serve via nginx
```

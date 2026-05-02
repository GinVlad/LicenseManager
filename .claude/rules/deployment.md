# Deployment Rules - LicenseManager

## Development
```bash
cp .env.example .env     # edit JWT_SECRET, ADMIN_PASSWORD
docker compose up -d     # postgres + backend + frontend
open http://localhost:5173
```

## Production Build
```bash
# Build backend image
docker build -t licensemanager .

# Build frontend static files
cd frontend && npm run build
# Serve dist/ behind nginx or upload to Vercel/Netlify
```

## Environment Variables (required in production)
```
PORT=8080
ENV=production
DATABASE_URL=postgres://user:pass@host:5432/licensemanager
JWT_SECRET=<32+ random bytes, keep secret>
ADMIN_EMAIL=you@example.com
ADMIN_PASSWORD=<strong password>
```

## Recommended Hosting
| Service | What |
|---------|------|
| Railway / Render | Backend container |
| Neon / Supabase | Managed PostgreSQL |
| Vercel / Netlify | Frontend static |

## Checklist Before Deploy
- [ ] JWT_SECRET is at least 32 random chars
- [ ] ADMIN_PASSWORD is strong
- [ ] DATABASE_URL points to production DB
- [ ] ENV=production (disables debug logs)
- [ ] CORS origins set to your frontend domain
- [ ] HTTPS enabled (HSTS header is already set)

## Railway Quick Deploy
```bash
# Install Railway CLI
npm install -g @railway/cli
railway login
railway init
railway add --database postgresql
railway up
```

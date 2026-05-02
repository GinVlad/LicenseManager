# Skill: Deploy to Railway

## Usage
`/deploy-railway`

## Steps

### 1. Install Railway CLI
```bash
npm install -g @railway/cli
railway login
```

### 2. Init Project
```bash
cd /Users/cheo/LicenseManager
railway init
```

### 3. Add PostgreSQL
```bash
railway add --database postgresql
# Railway sets DATABASE_URL automatically
```

### 4. Set Environment Variables
```bash
railway variables set JWT_SECRET=$(openssl rand -hex 32)
railway variables set ADMIN_EMAIL=you@example.com
railway variables set ADMIN_PASSWORD=yourpassword
railway variables set ENV=production
```

### 5. Deploy
```bash
railway up
```

### 6. Deploy Frontend
```bash
cd frontend && npm run build
# Upload dist/ to Vercel:
npx vercel dist/
# Set VITE_API_URL to your Railway backend URL
```

### 7. Update CORS
```bash
railway variables set CORS_ORIGINS=https://your-frontend.vercel.app
```

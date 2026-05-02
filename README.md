# LicenseManager

A self-hosted license management SaaS for your desktop apps.

- Multi-app: one server manages licenses for any number of apps
- Multi-HWID: one license can bind up to N devices (configurable per app)
- Expiry: set duration in days when issuing a key
- Revoke/extend: admin can revoke or extend any license at any time
- Admin dashboard: web UI to manage apps, issue/revoke licenses, view devices
- Go SDK: drop-in client for any Go/Wails app to validate licenses

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go + Gin |
| Database | PostgreSQL |
| Auth | JWT (admin) + API key (apps) |
| Frontend | React + Tailwind |

## Quick Start (Development)

```bash
# 1. Copy env
cp .env.example .env
# Edit .env — set ADMIN_PASSWORD and JWT_SECRET at minimum

# 2. Start everything
docker compose up -d

# 3. Open admin dashboard
open http://localhost:5173
# Login with ADMIN_EMAIL / ADMIN_PASSWORD from .env
```

## Production Deployment

```bash
# Build the image
docker build -t licensemanager .

# Run with your production .env
docker run -p 8080:8080 --env-file .env licensemanager
```

Serve the frontend (`npm run build`) behind nginx or deploy to Vercel/Netlify.

## Adding a New App

1. Open admin dashboard → **Apps** → **New App**
2. Fill in Name, Slug (e.g. `ebay-creator`), Max HWIDs per license
3. Copy the generated **API Key** — your desktop app needs it

## Issuing a License

1. Open **Licenses** → **Issue License**
2. Select the app, plan, thread limit, duration
3. The generated key (e.g. `EBAY-A1B2-C3D4-E5F6`) goes to your customer

## Integrating into a Go/Wails App

```go
import "github.com/cheo/licensemanager/sdk/go/license"

client := license.New(license.Config{
    ServerURL:  "https://licenses.yourdomain.com",
    AppSlug:    "ebay-creator",
    AppAPIKey:  "<api-key-from-admin-panel>",
    DeviceName: "Home PC",
})

info, err := client.Validate("EBAY-A1B2-C3D4-E5F6")
if err != nil {
    // show error to user: invalid key, expired, max devices, etc.
    log.Fatal(err)
}
// info.MaxThreads, info.ExpiresAt, info.HWID
```

## License Validation Flow

```
Client sends: { appSlug, key, hwid, deviceName }

Server checks (in order):
  1. App API key valid?          u2192 NO  u2192 401 Unauthorized
  2. License key exists for app? u2192 NO  u2192 { valid: false, error: "invalid key" }
  3. License is active?          u2192 NO  u2192 { valid: false, error: "license revoked" }
  4. Not expired?                u2192 NO  u2192 { valid: false, error: "license expired" }
  5. HWID already registered?    u2192 YES u2192 update last_seen, return valid
  6. HWID slots available?       u2192 NO  u2192 { valid: false, error: "max devices", devices: [...] }
  7. Register new HWID           u2192     u2192 { valid: true, expiresAt, maxThreads }
```

## API Reference

### Public (requires X-App-Key header)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/license/validate` | Validate key + register HWID |
| DELETE | `/api/v1/license/hwid?licenseId=...&hwidId=...` | Remove device |

### Admin (requires Authorization: Bearer <token>)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/admin/login` | Get JWT token |
| GET | `/api/v1/admin/apps` | List apps |
| POST | `/api/v1/admin/apps` | Create app |
| GET | `/api/v1/admin/licenses` | List licenses |
| POST | `/api/v1/admin/licenses` | Issue license |
| PUT | `/api/v1/admin/licenses/:id` | Update (revoke/extend/threads) |
| DELETE | `/api/v1/admin/licenses/:id` | Delete license |
| GET | `/api/v1/admin/licenses/:id/hwids` | List devices |
| DELETE | `/api/v1/admin/hwids/:id` | Remove device |

## Using This as a Template for a New App

1. Clone/copy this repo
2. Edit key and user in .env file, uncomment production if running on VPS
3. In admin panel → create an App with your app's slug
4. Copy the API key into your desktop app's config
5. In your Go app, import `sdk/go/license` and call `client.Validate(key)`
6. Block the app completely if `err != nil`

The SDK + server are app-agnostic. You can manage licenses for
10 different products from one LicenseManager instance.

# Database Agent - LicenseManager

## Role
PostgreSQL schema, migrations, and query optimization.

## Key Tables
```
admins(id, email, password_hash, role, created_at)
users(id, email, password_hash, stripe_customer_id, created_at)
applications(id, name, slug, api_key, max_hwid_per_license, created_at)
licenses(id, key, app_id, user_id, plan, max_threads, expires_at, is_active, note, created_at, last_validated_at)
license_hwids(id, license_id, hwid, device_name, first_seen_at, last_seen_at)
```

## Migration Process
1. Create `backend/internal/database/migrations/00N_description.up.sql`
2. Create `00N_description.down.sql` (always provide rollback)
3. Wrap in `BEGIN; ... COMMIT;`
4. Restart server - migrations run automatically

## Query Performance
- Validate lookup: `licenses.key` is indexed - fast
- HWID count: `license_hwids.license_id` is indexed - fast
- Admin list: add `WHERE app_id = $1` to filter by app - use `idx_licenses_app_id`
- For new WHERE columns: always add an index in the migration

## Common Queries
```sql
-- Validate (one query for license + app)
SELECT l.id, l.is_active, l.expires_at, l.max_threads, a.max_hwid_per_license
FROM licenses l
JOIN applications a ON l.app_id = a.id
WHERE l.key = $1 AND a.slug = $2;

-- HWID count
SELECT COUNT(*) FROM license_hwids WHERE license_id = $1;

-- Register HWID
INSERT INTO license_hwids (license_id, hwid, device_name)
VALUES ($1, $2, $3)
ON CONFLICT (license_id, hwid) DO UPDATE SET last_seen_at = NOW();
```

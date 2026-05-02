# Database Rules - LicenseManager

## Engine: PostgreSQL 16

## Schema Summary
```sql
admins          -- Operator accounts (you). Email + bcrypt hash.
users           -- Customer accounts (optional self-service).
applications    -- One row per product. Has slug + auto-generated API key.
licenses        -- One row per issued key. Belongs to an app.
license_hwids   -- One row per device bound to a license.
```

## Key Constraints 
- `applications.slug` UNIQUE  used in validate requests
- `applications.api_key` UNIQUE, DEFAULT `encode(gen_random_bytes(32), 'hex')`
- `licenses.key` UNIQUE  the key customers type
- `license_hwids(license_id, hwid)` UNIQUE  prevents duplicate bindings

## Migrations
- Location: `backend/internal/database/migrations/`
- Format: `00N_description.up.sql` / `.down.sql`
- Always wrap in `BEGIN; ... COMMIT;`
- Run automatically on server start via `database.RunMigrations()`
- Current: `001_init.up.sql` u2014 creates all tables

## Indexes
```sql
idx_licenses_key        ON licenses(key)           -- validate lookup
idx_licenses_app_id     ON licenses(app_id)        -- filter by app
idx_licenses_user_id    ON licenses(user_id)       -- user portal
idx_hwids_license_id    ON license_hwids(license_id) -- HWID count check
```

## Query Patterns
```go
// Single row
err := pool.QueryRow(ctx, "SELECT ... WHERE id = $1", id).Scan(&...)
if err == pgx.ErrNoRows { /* not found */ }

// Multiple rows
rows, err := pool.Query(ctx, "SELECT ... WHERE app_id = $1", appID)
defer rows.Close()
for rows.Next() { rows.Scan(&...) }

// Write
_, err := pool.Exec(ctx, "UPDATE ... SET x = $1 WHERE id = $2", x, id)

// Write + return
err := pool.QueryRow(ctx, "INSERT ... RETURNING id, ...").Scan(&id, &...)
```

## Adding a Column
1. Create `00N_add_column.up.sql`: `ALTER TABLE t ADD COLUMN c TYPE DEFAULT x;`
2. Create matching `down.sql`: `ALTER TABLE t DROP COLUMN c;`
3. Update struct in `models/models.go`
4. Update any affected Scan() calls

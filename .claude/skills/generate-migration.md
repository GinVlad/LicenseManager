# Skill: Generate Migration

## Usage
`/generate-migration [description]`

## Example
`/generate-migration add stripe_subscription_id to licenses`

## What It Does
Creates two files:
1. `backend/internal/database/migrations/00N_[description].up.sql`
2. `backend/internal/database/migrations/00N_[description].down.sql`

And tells you which struct in `models/models.go` to update.

## Template

### up.sql
```sql
BEGIN;

ALTER TABLE [table] ADD COLUMN [column] [type] DEFAULT [value];
CREATE INDEX IF NOT EXISTS idx_[table]_[column] ON [table]([column]);

COMMIT;
```

### down.sql
```sql
BEGIN;
DROP INDEX IF EXISTS idx_[table]_[column];
ALTER TABLE [table] DROP COLUMN [column];
COMMIT;
```

## Notes
- Migrations run automatically on server start
- Always wrap in BEGIN/COMMIT
- Always provide rollback
- Number format: 002, 003, ... (next after current highest)

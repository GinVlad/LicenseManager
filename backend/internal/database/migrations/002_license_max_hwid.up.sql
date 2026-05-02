BEGIN;

ALTER TABLE licenses ADD COLUMN max_hwid INTEGER NOT NULL DEFAULT 1;

-- Backfill existing licenses from their app's default
UPDATE licenses l
SET max_hwid = a.max_hwid_per_license
FROM applications a
WHERE l.app_id = a.id;

COMMIT;

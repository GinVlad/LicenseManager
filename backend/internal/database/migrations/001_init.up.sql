BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Admins (you, the operator)
CREATE TABLE IF NOT EXISTS admins (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email        VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role         VARCHAR(20) DEFAULT 'admin',
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

-- Users (your customers, optional self-service portal)
CREATE TABLE IF NOT EXISTS users (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email              VARCHAR(255) UNIQUE NOT NULL,
    password_hash      VARCHAR(255),
    stripe_customer_id VARCHAR(50),
    created_at         TIMESTAMPTZ DEFAULT NOW()
);

-- Applications (one row per app you sell: "eBay Creator", "Amazon Tool", etc.)
CREATE TABLE IF NOT EXISTS applications (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                  VARCHAR(100) NOT NULL,
    slug                  VARCHAR(50) UNIQUE NOT NULL,
    api_key               VARCHAR(64) UNIQUE NOT NULL DEFAULT encode(gen_random_bytes(32), 'hex'),
    max_hwid_per_license  INT DEFAULT 2,
    created_at            TIMESTAMPTZ DEFAULT NOW()
);

-- Licenses
CREATE TABLE IF NOT EXISTS licenses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key               VARCHAR(32) UNIQUE NOT NULL,
    app_id            UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    user_id           UUID REFERENCES users(id) ON DELETE SET NULL,
    plan              VARCHAR(20) NOT NULL DEFAULT 'basic',
    max_threads       INT DEFAULT 5,
    expires_at        TIMESTAMPTZ NOT NULL,
    is_active         BOOLEAN DEFAULT true,
    note              TEXT DEFAULT '',
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    last_validated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_licenses_key ON licenses(key);
CREATE INDEX IF NOT EXISTS idx_licenses_app_id ON licenses(app_id);
CREATE INDEX IF NOT EXISTS idx_licenses_user_id ON licenses(user_id);

-- HWID bindings (multi-device per license)
CREATE TABLE IF NOT EXISTS license_hwids (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_id   UUID NOT NULL REFERENCES licenses(id) ON DELETE CASCADE,
    hwid         VARCHAR(64) NOT NULL,
    device_name  VARCHAR(100) DEFAULT '',
    first_seen_at TIMESTAMPTZ DEFAULT NOW(),
    last_seen_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(license_id, hwid)
);

CREATE INDEX IF NOT EXISTS idx_hwids_license_id ON license_hwids(license_id);

COMMIT;

CREATE TABLE IF NOT EXISTS
    credentials (
        "user_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        "password_hash" BYTEA NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS
    idx_credentials_updated_at ON credentials ("updated_at");
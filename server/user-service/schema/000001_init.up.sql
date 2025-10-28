CREATE TABLE IF NOT EXISTS
    users (
        "user_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        "username" TEXT NOT NULL UNIQUE,
        "email" TEXT NOT NULL UNIQUE,
        "name" TEXT NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
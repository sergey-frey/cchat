CREATE TABLE IF NOT EXISTS
    chats (
        "chat_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        "name" TEXT NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    user_chats (
        "chat_id" UUID NOT NULL,
        "user_id" UUID NOT NULL
    );
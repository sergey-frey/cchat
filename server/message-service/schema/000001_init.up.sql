CREATE TABLE IF NOT EXISTS
    messages (
        "message_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        "type" TEXT NOT NULL,
        "content" TEXT NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        "chat_id" UUID NOT NULL,
        "author_id" UUID NOT NULL
    );
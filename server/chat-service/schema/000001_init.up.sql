CREATE TABLE IF NOT EXISTS
    chats (
        "chat_id" UUID PRIMARY KEY DEFAULT gen_random_uuid()
    );

CREATE TABLE IF NOT EXISTS
    user_chats (
        "chat_id" UUID NOT NULL,
        "user_id" UUID NOT NULL
    );
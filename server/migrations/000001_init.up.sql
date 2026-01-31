CREATE TABLE IF NOT EXISTS
    chats (
        "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    );

CREATE TABLE IF NOT EXISTS
    users (
        "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        "username" TEXT NOT NULL UNIQUE,
        "email" TEXT NOT NULL UNIQUE,
        "pass_hash" BYTEA NOT NULL,
        "name" TEXT NOT NULL
    );

CREATE TABLE IF NOT EXISTS
    messages (
        "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        "type" TEXT NOT NULL,
        "content" TEXT NOT NULL,
        "date" DATE NOT NULL,
        "chat_id" INT NOT NULL,
        "author_id" INT NOT NULL,
        CONSTRAINT "messages_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES chats ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
        CONSTRAINT "messages_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES users ("id") ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    user_chats (
        "chat_id" BIGINT NOT NULL,
        "user_id" BIGINT NOT NULL,
        CONSTRAINT "user_chats_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES chats ("id") ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT "user_chats_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES users ("id") ON DELETE CASCADE ON UPDATE CASCADE
    );
CREATE TABLE chats (
    "id" SERIAL PRIMARY KEY
);

CREATE TABLE users (
    "id" SERIAL PRIMARY KEY,
    "username" TEXT NOT NULL UNIQUE,
    "email" TEXT NOT NULL UNIQUE,
    "pass_hash" BYTEA NOT NULL
    -- "name" TEXT NOT NULL 
    --"socket" TEXT
);

CREATE TABLE messages (
    "id" SERIAL PRIMARY KEY,
    "type" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "date" DATE NOT NULL,
    "chatId" INT NOT NULL,
    "authorId" INT NOT NULL,
    CONSTRAINT "Message_chatId_fkey" FOREIGN KEY ("chatId") REFERENCES chats ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT "Message_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES users ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE _ChatToUser (
    "A" INT NOT NULL,
    "B" INT NOT NULL,
    CONSTRAINT "_ChatToUser_A_fkey" FOREIGN KEY ("A") REFERENCES chats ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "_ChatToUser_B_fkey" FOREIGN KEY ("B") REFERENCES users ("id") ON DELETE CASCADE ON UPDATE CASCADE
);
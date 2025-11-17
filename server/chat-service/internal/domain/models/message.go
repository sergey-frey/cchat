package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	UUID          uuid.UUID     `json:"message_id"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	ChatID      uuid.UUID       `json:"chat_id"`
	AuthorID    uuid.UUID `json:"author_id"`
}

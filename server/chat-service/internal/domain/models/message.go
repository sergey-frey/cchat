package models

import "time"

type Message struct {
	ID          int64     `json:"message_id"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	ChatID      int       `json:"chat_id"`
	AuthorID    int       `json:"author_id"`
}

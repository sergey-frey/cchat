package models

import "time"

type Message struct {
	ID int64
	TypeOfMessage string
	Content string
	Date time.Time
	ChatID int
	AuthorID int
}
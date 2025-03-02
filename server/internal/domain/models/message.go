package models

import "time"

type Message struct {
	ID int64
	TypeMessage string
	Content string
	Date time.Time
	HasBeenEdited bool
	// Chat
	// User
}
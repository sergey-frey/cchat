package models

import "github.com/google/uuid"

type Cursor struct {
	NextCursor uuid.UUID `json:"next_cursor"`
	PrevCursor uuid.UUID `json:"prev_cursor"`
}

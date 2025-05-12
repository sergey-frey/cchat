package models

type Cursor struct {
	NextCursor int64 `json:"next_cursor"`
	PrevCursor int64 `json:"prev_cursor"`
}
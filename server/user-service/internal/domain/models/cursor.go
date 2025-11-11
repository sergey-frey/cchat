package models

type Cursor struct {
	NextCursor string `json:"next_cursor"`
	HasNextPage bool `json:"has_next_page"`
}

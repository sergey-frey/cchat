package models

type NewChat struct {
	Users []int64 `json:"users"`
}

type Chat struct {
	ID          int64 `json:"id"`
	Users       []UserInfo
	LastMessage *Message `json:"last_message"`
}

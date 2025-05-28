package models

type NewChat struct {
	Users []int64 `json:"users"`
}

type Chat struct {
	ID int64
	Users []UserInfo
	Messages []Message
}
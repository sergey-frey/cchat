package models

import "github.com/google/uuid"

type NewChat struct {
	ChatName string `json:"chat_name"`
	Users []uuid.UUID `json:"users"`
}

type Chat struct {
	UUID          uuid.UUID `json:"id"`
	Users       []UserInfo
	LastMessage *Message `json:"last_message"`
}

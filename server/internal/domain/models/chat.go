package models

type Chat struct {
	ID int64
	Users []User
	Messages []Message
	Secret string
}
package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUsersNotFound = errors.New("users not found")
	ErrUserExists = errors.New("user already exists")
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists = errors.New("email already exists")
	ErrFailedToCreateChat = errors.New("failed to create new chat")
	ErrFailedToAddUsersInChat = errors.New("failed to add users in chat")
)
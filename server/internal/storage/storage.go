package storage

import "errors"

var (
	ErrUserExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists = errors.New("email already exists")
)
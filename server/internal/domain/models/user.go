package models

type User struct {
	ID int64
	Username string `json:"username"`
	Email string `json:"email" validate:"required"`
	PassHash []byte `json:"password" validate:"required"`
}

type RegisterUser struct {
	Username string `json:"username" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUser struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type NormalizedUser struct {
	ID int64
	Username string
	Email string
}
package models

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email" validate:"required"`
	PassHash []byte    `json:"password" validate:"required"`
}

type RegisterUser struct {
	Email    string `json:"email" validate:"required,email" example:"example@mail.com"`
	Password string `json:"password" validate:"required,gte=8" example:"12345678"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email" example:"example@mail.com"`
	Password string `json:"password" validate:"required" example:"12345678"`
}

type UserInfo struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
}

type NewPassword struct {
	Email            string `json:"email" validate:"required,email" example:"example@mail.com"`
	PreviousPassword string `json:"previous_password" validate:"required" example:"12345678"`
	NewPassword      string `json:"new_password" validate:"required,gte=8" example:"123456789"`
}

type NormalizedUser struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func UserToNormalized(user *User) NormalizedUser {
	return NormalizedUser{
		UUID:     user.UUID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func InfoToNormalized(info *UserInfo) NormalizedUser {
	return NormalizedUser{
		UUID:     info.UUID,
		Username: info.Username,
		Email:    info.Email,
	}
}

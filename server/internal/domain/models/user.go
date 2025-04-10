package models

type User struct {
	ID       int64
	Username string `json:"username"`
	Email    string `json:"email" validate:"required"`
	PassHash []byte `json:"password" validate:"required"`
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
	ID       int64
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type NewUserInfo struct {
	PreviousPassword *string `json:"previous_password,omitempty"`
	NewPassword      *string `json:"new_password,omitempty" validate:"omitempty,gte=8"`
	Email            *string `json:"email,omitempty"`
	Username         *string `json:"username,omitempty"`
	Name             *string `json:"name,omitempty" validate:"omitempty,gte=1"`
}

type NormalizedUser struct {
	ID       int64
	Username string
	Email    string
}

func UserToNormalized(user *User) NormalizedUser {
	return NormalizedUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func InfoToNormalized(info *UserInfo) NormalizedUser {
	return NormalizedUser{
		ID:       info.ID,
		Username: info.Username,
		Email:    info.Email,
	}
}

package models

type User struct {
	ID       int64
	Username string `json:"username"`
	Email    string `json:"email" validate:"required"`
	PassHash []byte `json:"password" validate:"required"`
}

type RegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserInfo struct {
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Name     string `json:"name" db:"name"`
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

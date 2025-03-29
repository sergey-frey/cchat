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
	ID       int64
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type NewUserInfo struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	Name        string `json:"name,omitempty"`
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

func InfoToNormalized(info *UserInfo) NormalizedUser{
	return NormalizedUser{
		ID:       info.ID,
		Username: info.Username,
		Email:    info.Email,
	}
}

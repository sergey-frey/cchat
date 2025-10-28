package models

type User struct {
	ID       int64  `json:"id"`
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
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type NewUserInfo struct {
	Email            string `json:"email,omitempty" validate:"omitempty,email" example:"example@mail.com"`
	PreviousPassword string `json:"previous_password,omitempty" validate:"omitempty" example:"12345678"`
	NewPassword      string `json:"new_password,omitempty" validate:"omitempty,gte=8" example:"123456789"`
	Username         string `json:"username,omitempty" validate:"omitempty" example:"arnold2004"`
	Name             string `json:"name,omitempty" validate:"omitempty,min=1" example:"Arnold"`
}

type NormalizedUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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

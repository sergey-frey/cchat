package tests

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/sergey-frey/cchat/internal/domain/models"
)

const (
	host = "localhost:8040"
)

func TestCchatAuth_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/cchat/auth/register").
		WithJSON(models.RegisterUser{
			Username: gofakeit.Username(),
			Email: gofakeit.Email(),
			Password: generateNewPassword(),
		}).
		Expect().
		Status(200)
}

func TestRegister_FailCases(t *testing.T) {
	testCases := []struct{
		name string
		username string
		email string
		password string
		expectedErr string
	}{
		{
			name: "Register without username",
			username: "",
			email: gofakeit.Email(),
			password: generateNewPassword(),
			expectedErr: "field Username is a required field",
		},
		{
			name: "Register without email",
			username: gofakeit.Username(),
			email: "",
			password: generateNewPassword(),
			expectedErr: "field Email is a required field",
		},
		{
			name: "Register without password",
			username: gofakeit.Username(),
			email: gofakeit.Email(),
			password: "",
			expectedErr: "field Password is a required field",
		},
		{
			name: "Register with incorrect email",
			username: gofakeit.Username(),
			email: gofakeit.Username(),
			password: generateNewPassword(),
			expectedErr: "field Email is not valid",
		},
		{
			name: "Register with incorrect password",
			username: gofakeit.Username(),
			email: gofakeit.Email(),
			password: generateInvalidPassword(),
			expectedErr: "field Password is not valid",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host: host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/cchat/auth/register").
				WithJSON(models.RegisterUser{
					Username: tt.username,
					Email: tt.email,
					Password: tt.password,
			}).Expect().JSON().Object()

			if tt.expectedErr != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.expectedErr)

				return
			}

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	testCases := []struct{
		name string
		username string
		email string
		password string
		expectedErr string
	}{
		{
			name: "Login without email",
			email: gofakeit.Email(),
			password: generateInvalidPassword(),
			expectedErr: "email is required",
		},
		{
			name: "Login without password",
			email: gofakeit.Email(),
			password: "",
			expectedErr: "password is required",
		},
		{
			name: "Login without incorrect email",
			email: gofakeit.Email(),
			password: generateNewPassword(),
			expectedErr: "email is incorrect",
		},
		{
			name: "Login without incorrect password",
			email: gofakeit.Email(),
			password: generateInvalidPassword(),
			expectedErr: "password is incorrect",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host: host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/cchat/auth/login").
				WithJSON(models.RegisterUser{
					Email: tt.email,
					Password: tt.password,
			}).Expect().JSON().Object()

			if tt.expectedErr != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.expectedErr)

				return
			}

		})
	}
}

func generateNewPassword() string {
	return gofakeit.Password(true, true, true, false, false, 20)
}

func generateInvalidPassword() string {
	return gofakeit.Password(true, false, true, false, false, 6)
}
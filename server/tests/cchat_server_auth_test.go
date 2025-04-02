package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/sergey-frey/cchat/internal/domain/models"
)

const (
	host = "localhost:8040"
	normalLengthPass = 12
	notEnoughLengthPass = 6
)

func TestCchatAuth_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	email := gofakeit.Email()
	password := randomFakePassword(normalLengthPass)

	e.POST("/cchat/auth/register").
		WithJSON(models.RegisterUser{
			Email: email,
			Password: password,
		}).
		Expect().
		Status(http.StatusOK)

	e.POST("/cchat/auth/login").
		WithJSON(models.LoginUser{
			Email: email,
			Password: password,
		}).
		Expect().
		Status(http.StatusOK)
}

func TestRegister_FailCases(t *testing.T) {
	cases := []struct {
		name string
		email string
		password string
		respError string
	}{
		{
			name: "Register with empty email",
			email: "",
			password: randomFakePassword(normalLengthPass),
			respError: "field Email is a required field",
		},
		{
			name: "Register with empty password",
			email: gofakeit.Email(),
			password: "",
			respError: "field Password is a required field",
		},
		{
			name: "Register with both empty",
			email: "",
			password: "",
			respError: "field Email is a required field, field Password is a required field",
		},
		{
			name: "Register with invalid email",
			email: gofakeit.Username(),
			password: randomFakePassword(normalLengthPass),
			respError: "field Email must be of the email type",
		},
		{
			name: "Register with invalid passsword",
			email: gofakeit.Email(),
			password: randomFakePassword(notEnoughLengthPass),
			respError: "field Password must have at least 8 characters",
		},
	}

	for _, tt := range cases {

		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host: "localhost:8040",
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/cchat/auth/register").
				WithJSON(models.RegisterUser{
					Email: tt.email,
					Password: tt.password,
				}).Expect().JSON().Object()

			if tt.respError != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.respError)

				return
			}
		})
	}
}
		


func TestLogin_FailCases(t *testing.T) {
	cases := []struct{
		name string
		username string
		email string
		password string
		expectedErr string
	}{
		{
			name: "Login without email",
			email: "",
			password: randomFakePassword(normalLengthPass),
			expectedErr: "field Email is a required field",
		},
		{
			name: "Login without password",
			email: gofakeit.Email(),
			password: "",
			expectedErr: "field Password is a required field",
		},
		{
			name: "Login with invalid email",
			email: gofakeit.Username(),
			password: randomFakePassword(normalLengthPass),
			expectedErr: "field Email must be of the email type",
		},
		{
			name: "Login with invalid password",
			email: gofakeit.Email(),
			password: randomFakePassword(notEnoughLengthPass),
			expectedErr: "invalid email or password",
		},
	}

	for _, tt := range cases {
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

func randomFakePassword(length int) string {
	return gofakeit.Password(true, true, true, false, false, length)
}
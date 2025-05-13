package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/sergey-frey/cchat/internal/domain/models"
)

func TestGetUserInfo_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	email := gofakeit.Email()
	password := randomFakePassword(normalLengthPass)

	e.POST("/cchat/auth/register").
		WithJSON(models.RegisterUser{
			Email:    email,
			Password: password,
		}).
		Expect().
		Status(http.StatusOK)

	e.GET("/cchat/user/profile/{username}").
	WithPath("username", "weeq").
	Expect().Status(http.StatusOK)
}

func TestUpdateUserInfo_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	email := gofakeit.Email()
	previousPassword := randomFakePassword(normalLengthPass)
	newPassword := randomFakePassword(normalLengthPass)
	username := gofakeit.Username()
	name := gofakeit.Name()

	e.POST("/cchat/auth/register").
		WithJSON(models.RegisterUser{
			Email:    email,
			Password: previousPassword,
		}).
		Expect().
		Status(http.StatusOK)

	e.PATCH("/cchat/user/update").
		WithJSON(models.NewUserInfo{
			PreviousPassword: previousPassword,
			NewPassword:      newPassword,
			Username:         username,
			Name:             name,
		}).
		Expect().
		Status(http.StatusOK)
}

// func TestGetUserInfo_FailCases(t *testing.T) {
// 	cases := []struct {
// 		name string
// 	}{
// 		{},
// 	}
// }

func TestUpdatePassword_FailCases(t *testing.T) {
	email := gofakeit.Email()
	password := randomFakePassword(normalLengthPass)

	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	resp := e.POST("/cchat/auth/register").
		WithJSON(models.RegisterUser{
			Email:    email,
			Password: password,
		}).Expect().JSON().Object()

	cases := []struct {
		name             string
		email            string
		previousPassword string
		newPassword      string
		expectedErr      string
	}{
		{
			name:             "Update password with incorrect Previous password",
			email:            email,
			previousPassword: randomFakePassword(normalLengthPass),
			newPassword:      randomFakePassword(normalLengthPass),
			expectedErr:      "passwords don't match",
		},
		{
			name:             "Update password with invalid New password",
			email:            email,
			previousPassword: password,
			newPassword:      randomFakePassword(notEnoughLengthPass),
			expectedErr:      "field NewPassword must have at least 8 characters",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp = e.PATCH("/cchat/user/update").
				WithJSON(models.NewUserInfo{
					PreviousPassword: tt.previousPassword,
					NewPassword:      tt.newPassword,
				}).Expect().JSON().Object()

			if tt.expectedErr != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.expectedErr)

				return
			}

		})
	}
}

func TestUpdateUsername_FailCases(t *testing.T) {
	cases := []struct {
		name             string
		email            string
		previousPassword string
		newPassword      string
		username         string
		personName       string
		expectedErr      string
	}{
		{
			name:             "Update username with already existed username",
			email:            gofakeit.Email(),
			previousPassword: randomFakePassword(normalLengthPass),
			username:         "stepan42k",
			expectedErr:      "username already exists",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			e.POST("/cchat/auth/register").
				WithJSON(models.RegisterUser{
					Email:    tt.email,
					Password: tt.previousPassword,
				}).Expect().JSON().Object()

			resp := e.PATCH("/cchat/user/update").
				WithJSON(models.NewUserInfo{
					Username:         tt.username,
					Name:             tt.name,
				}).Expect().JSON().Object()

			if tt.expectedErr != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.expectedErr)

				return
			}

		})
	}

}

func TestUpdateName_FailCases(t *testing.T) {
	cases := []struct {
		name             string
		email            string
		previousPassword string
		newPassword      string
		username         string
		personName       string
		expectedErr      string
	}{
		{
			name:             "Update name with invalid name",
			email:            gofakeit.Email(),
			previousPassword: randomFakePassword(normalLengthPass),
			personName:       "",
			expectedErr:      "field Name must have at least 1 characters",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			e.POST("/cchat/auth/register").
				WithJSON(models.RegisterUser{
					Email:    tt.email,
					Password: tt.previousPassword,
				}).Expect().JSON().Object()

			resp := e.PATCH("/cchat/user/update").
				WithJSON(models.NewUserInfo{
					Name: tt.personName,
				}).Expect().JSON().Object()

			if tt.expectedErr != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.expectedErr)

				return
			}

		})
	}
}

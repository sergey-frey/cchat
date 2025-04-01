package cookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
)


func TakeUserInfo(w http.ResponseWriter, r *http.Request) (string, error) {
	user, err := CheckCookie(w, r)
	if err != nil {
		return "", fmt.Errorf("error with taking cookie")
	}

	return user.Username, nil
}


func SetCookie(w http.ResponseWriter, accessToken string, refreshToken string) {
	cookie1 := &http.Cookie{
		Name: "access_token",
		Value: accessToken,
		Path: "/",
		HttpOnly: true,
		Secure: false,
	}
	http.SetCookie(w, cookie1)

	cookie2 := &http.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		Path: "/",
		HttpOnly: true,
		Secure: false,
	}
	http.SetCookie(w, cookie2)
}


func CheckCookie(w http.ResponseWriter, r *http.Request) (*models.NormalizedUser, error) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		return HandlerError(err)
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return HandlerError(err)
	}

	accessToken, refreshToken := accessCookie.Value, refreshCookie.Value

	accessToken, refreshToken, user, err := jwt.VerifyAccessToken(accessToken, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error with token: %w", err)
	}

	if accessToken != "" {
		SetCookie(w, accessToken, refreshToken)
	}
	
	return user, nil
}


func DeleteCookie(w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name: "access_token",
		Value: "",
		Path: "/",
		HttpOnly: true,
		Secure: false,
		Expires: time.Now(),
	}
	http.SetCookie(w, cookie)

	cookie = &http.Cookie{
		Name: "refresh_token",
		Value: "",
		Path: "/",
		HttpOnly: true,
		Secure: false,
		Expires: time.Now(),
	}
	http.SetCookie(w, cookie)
}

func HandlerError(err error) (*models.NormalizedUser, error) {
	switch err {
		case http.ErrNoCookie:
			return nil, http.ErrNoCookie
		default:
			return nil, fmt.Errorf("server error")
		}
}
package cookie

import (
	"fmt"
	"net/http"
	"time"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
)


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


func CheckCookie(w http.ResponseWriter, r *http.Request) (bool, error) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		return HandlerError(err)
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return HandlerError(err)
	}

	accessToken, refreshToken := accessCookie.Value, refreshCookie.Value

	accessToken, refreshToken, flag, err := jwt.VerifyAccessToken(accessToken, refreshToken)
	if err != nil {
		return false, fmt.Errorf("error with token: %w", err)
	}

	if accessToken != "" {
		SetCookie(w, accessToken, refreshToken)
	}
	
	return flag, nil
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

func HandlerError(err error) (bool, error) {
	switch err {
		case http.ErrNoCookie:
			return false, http.ErrNoCookie
		default:
			return false, fmt.Errorf("server error")
		}
}
package cookie

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sergey-frey/cchat/internal/lib/jwt"
)

func SetCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name: "jwt_token",
		Value: token,
		Path: "/",
		HttpOnly: true,
		Secure: false,
		Expires: time.Now().Add(2 * time.Minute),
	}

	http.SetCookie(w, cookie)
}

func CheckCookie(r *http.Request) error {
	cookie, err := r.Cookie("jwt_token")

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return fmt.Errorf("cookie not found")
		default:
			return fmt.Errorf("server error")
		}
	}

	tokenString := cookie.Value

	jwt.VerifyToken(tokenString)

	return nil
}

func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name: "jwt_token",
		Value: "",
		Expires: time.Now(),
	}

	http.SetCookie(w, cookie)
}
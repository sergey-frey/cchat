package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sergey-frey/cchat/internal/domain/models"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["secret"] = "somesecret"
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte("somesecret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "somesecret", nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
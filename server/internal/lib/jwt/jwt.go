package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sergey-frey/cchat/internal/domain/models"
)


const (
	accessDuration = 10 * time.Second
	refreshDuration = 30 * time.Second
)

var (
	ErrUserUnauthorized = errors.New("user unauthorized")
)


func NewPairTokens(user models.User) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(accessDuration).Unix()

	accessTokenString, err := accessToken.SignedString([]byte("somesecret"))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(refreshDuration).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil

}


func VerifyAccessToken(accessToken string, refreshToken string) (string, string, bool, error) {
	accesstoken, _ := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte("somesecret"), nil
	})

	if !accesstoken.Valid {
		claims := accesstoken.Claims.(jwt.MapClaims)
		newAccessToken, newRefreshToken, flag, err := VerifyRefreshToken(claims, refreshToken)	
		return newAccessToken, newRefreshToken, flag, err
	}

	return "", "" , true, nil 
}


func VerifyRefreshToken(claims jwt.MapClaims, refreshToken string) (string, string, bool, error) {

	refreshtoken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
			return []byte("secret"), nil
		})
	
	if err != nil {
		return "", "", false, err
	}

	if !refreshtoken.Valid {
		return "", "", false, ErrUserUnauthorized
	}

	var user = models.User{
		ID: int64(claims["uid"].(float64)),
		Email: claims["email"].(string),
	}

	newAccessToken, newRefreshToken, err := NewPairTokens(user)
	if err != nil {				
		return "", "", false, err
	}

	return newAccessToken, newRefreshToken, true, nil
}
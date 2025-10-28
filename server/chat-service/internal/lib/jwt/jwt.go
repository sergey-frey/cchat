package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sergey-frey/cchat/server/chat-service/internal/domain/models"
)

const (
	accessDuration  = 15 * time.Minute
	refreshDuration = 43200 * time.Minute
)

var (
	ErrUserUnauthorized = errors.New("user unauthorized")
)

func NewPairTokens(user models.NormalizedUser) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["username"] = user.Username
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

func VerifyAccessToken(accessToken string, refreshToken string) (string, string, *models.NormalizedUser, error) {
	accesstoken, _ := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte("somesecret"), nil
	})

	claims := accesstoken.Claims.(jwt.MapClaims)
	var user = models.NormalizedUser{
		ID:       int64(claims["uid"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}

	if !accesstoken.Valid {
		newAccessToken, newRefreshToken, user, err := VerifyRefreshToken(user, refreshToken)
		return newAccessToken, newRefreshToken, user, err
	}

	return "", "", &user, nil
}

func VerifyRefreshToken(user models.NormalizedUser, refreshToken string) (string, string, *models.NormalizedUser, error) {

	refreshtoken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return "", "", nil, err
	}

	if !refreshtoken.Valid {
		return "", "", nil, ErrUserUnauthorized
	}

	newAccessToken, newRefreshToken, err := NewPairTokens(user)
	if err != nil {
		return "", "", nil, err
	}

	return newAccessToken, newRefreshToken, &user, nil
}

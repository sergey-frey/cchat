package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sergey-frey/cchat/server/auth-service/internal/domain/models"
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
	claims["uuid"] = user.UUID
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
	claims := &jwt.MapClaims{}

	accesstoken, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("somesecret"), nil
	})

	if err != nil || !accesstoken.Valid {

		user, userErr := extractUserFromClaims(*claims)
		if userErr != nil {
			return "", "", nil, ErrUserUnauthorized
		}
		
		return VerifyRefreshToken(*user, refreshToken)
	}

	user, userErr := extractUserFromClaims(*claims)
	if userErr != nil {
		return "", "", nil, userErr
	}

	return "", "", user, nil
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

func extractUserFromClaims(claims jwt.MapClaims) (*models.NormalizedUser, error) {
	uuidStr, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'uuid' in token claims")
	}
	
	userUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, errors.New("failed to parse 'uuid' from token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'username' in token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'email' in token claims")
	}

	return &models.NormalizedUser{
		UUID:     userUUID,
		Username: username,
		Email:    email,
	}, nil
}

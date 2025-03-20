package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
	genusername "github.com/sergey-frey/cchat/internal/lib/username"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	log *slog.Logger
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, email string, passHash []byte) (models.NormalizedUser, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)


func New(userSaver UserSaver, log *slog.Logger, userProvider UserProvider) *AuthService {
	return &AuthService{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log: log,
	}
}


func (a *AuthService) Login(ctx context.Context, loginUser models.LoginUser) (models.NormalizedUser, string, string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", loginUser.Email),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, loginUser.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))

		return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(loginUser.Password)); err != nil {
		log.Warn("invalid credentials", sl.Err(err))
		return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	normalUser := models.UserToNormalized(user)

	accessToken, refreshToken, err := jwt.NewPairTokens(normalUser)
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))

		return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, err)
	}

	return normalUser, accessToken, refreshToken, err
}


func (a *AuthService) RegisterNewUser(ctx context.Context, email string, pass string) (models.NormalizedUser, string, string, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, err)
	}

	username := genusername.GenerateUsername()

	user, err := a.usrSaver.SaveUser(ctx, username, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	accessToken, refreshToken, err := jwt.NewPairTokens(user)
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))

		return models.NormalizedUser{}, "", "", fmt.Errorf("%s: %w", op, err)
	}

	return user, accessToken, refreshToken, nil
}
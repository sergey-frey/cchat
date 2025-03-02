package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
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


func (a *AuthService) Login(ctx context.Context, email string, password string) (models.NormalizedUser, string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email), //("username", email)
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return models.NormalizedUser{}, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))

		return models.NormalizedUser{}, "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", sl.Err(err))
		return models.NormalizedUser{}, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, 2 * time.Minute)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))

		return models.NormalizedUser{}, "", fmt.Errorf("%s: %w", op, err)
	}

	normalUser := UserToNormalized(user)

	return normalUser, token, nil
}


func (a *AuthService) RegisterNewUser(ctx context.Context, username string, email string, pass string) (models.NormalizedUser, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return models.NormalizedUser{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.usrSaver.SaveUser(ctx, username, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return models.NormalizedUser{}, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return models.NormalizedUser{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return user, nil
}


func UserToNormalized(user models.User) models.NormalizedUser {
	return models.NormalizedUser{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
	}
}

package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sergey-frey/cchat/server/auth-service/internal/domain/models"
	"github.com/sergey-frey/cchat/server/auth-service/internal/lib/jwt"
	"github.com/sergey-frey/cchat/server/auth-service/internal/lib/logger/sl"
	storage "github.com/sergey-frey/cchat/server/auth-service/internal/provider"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider interface {
	Register(ctx context.Context, id uuid.UUID, passHash []byte) (error)
	Login(ctx context.Context, uuid uuid.UUID) ([]byte, error)
	Password(ctx context.Context, uuid uuid.UUID) ([]byte, error)
	ChangePassword(ctx context.Context, uid uuid.UUID, newPasswordHash []byte) error
	RefreshPassword(ctx context.Context) error
}

type UserProvider interface {
	GetUser(ctx context.Context, userID string) (*models.NormalizedUser, error)
	CreateUser(ctx context.Context, email string) (*models.NormalizedUser, error)
}

type AuthService struct {
	auth AuthProvider
	user UserProvider
	log  *slog.Logger
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedToGetUser    = errors.New("failed to get user")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

func New(auth AuthProvider, user UserProvider, log *slog.Logger) *AuthService {
	return &AuthService{
		auth: auth,
		user: user,
		log:  log,
	}
}

//go:generate go run github.com/vektra/mockery/v2@v2.53 --name=Auth
func (a *AuthService) Login(ctx context.Context, loginUser models.LoginUser) (*models.NormalizedUser, string, string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", loginUser.Email),
	)

	log.Info("attempting to login user")

	user, err := a.user.GetUser(ctx, loginUser.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	passHash, err := a.auth.Login(ctx, user.UUID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(passHash, []byte(loginUser.Password)); err != nil {
		log.Warn("invalid credentials", sl.Err(err))
		return nil, "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	accessToken, refreshToken, err := jwt.NewPairTokens(*user)
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	return user, accessToken, refreshToken, err
}

func (a *AuthService) RegisterNewUser(ctx context.Context, registerUser models.RegisterUser) (*models.NormalizedUser, string, string, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", registerUser.Email),
	)

	log.Info("registering user")

	user, err := a.user.CreateUser(ctx, registerUser.Email)
	if err != nil {
		log.Error("failed to get user", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, ErrFailedToGetUser)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.auth.Register(ctx, user.UUID, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	accessToken, refreshToken, err := jwt.NewPairTokens(*user)
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	return user, accessToken, refreshToken, nil
}

func (a *AuthService) ChangePassword(ctx context.Context, newPassword models.NewPassword) (*models.NormalizedUser, string, string, error) {
	const op = "auth.ChangePassword"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", newPassword.Email),
	)

	log.Info("registering user")

	user, err := a.user.GetUser(ctx, newPassword.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		log.Error("failed to get user", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}
	
	oldPasswordHash, err := a.auth.Password(ctx, user.UUID)
	if err != nil {
		log.Error("failed to get old password", sl.Err(err))
		
		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(oldPasswordHash, []byte(newPassword.PreviousPassword)); err != nil {
		log.Error("failed to compare passwords", sl.Err(ErrInvalidCredentials))
		
		return nil, "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(newPassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.auth.ChangePassword(ctx, user.UUID, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	accessToken, refreshToken, err := jwt.NewPairTokens(*user)
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	return user, accessToken, refreshToken, nil
}


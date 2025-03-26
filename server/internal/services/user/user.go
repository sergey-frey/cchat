package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceProvider interface {
	GetUser(ctx context.Context, username string) (info *models.UserInfo, err error)
	ChangeUsername(ctx context.Context, oldUsername string, newUsername string) (info *models.NormalizedUser, err error)
	ChangeName(ctx context.Context, username string, newName string) (info *models.UserInfo, err error)
	ChangePassword(ctx context.Context, username string, newPasswordHash []byte) (flag bool, err error)
	GetPassword(ctx context.Context, username string) (passHash []byte, err error)
}

type UserDataService struct {
	userServiceProvider UserServiceProvider
	log          *slog.Logger
}

func New(userProvider UserServiceProvider, log *slog.Logger) *UserDataService {
	return &UserDataService{
		userServiceProvider: userProvider,
		log:          log,
	}
}


func (u *UserDataService) GetUser(ctx context.Context, username string) (*models.UserInfo, error) {
	const op = "services.user.GetUser"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("getting user information")

	info, err := u.userServiceProvider.GetUser(ctx, username)
	if err != nil {
		log.Error("failed to get user", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got info")

	return info, nil
}


func (u *UserDataService) ChangeUsername(ctx context.Context, oldUsername string, newUsername string) (*models.NormalizedUser, string, string, error) {
	const op = "services.user.ChangeUsername"

	log := u.log.With(
		slog.String("op", op),
		slog.String("oldUsername", oldUsername),
		slog.String("newUsername", newUsername),
	)

	log.Info("changing username")

	user, err := u.userServiceProvider.ChangeUsername(ctx, oldUsername, newUsername)
	if err != nil {
		log.Error("failed to change username")

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, refreshToken, err := jwt.NewPairTokens(*user)
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))

		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("username changed")

	return user, accessToken, refreshToken, nil
}


func (u *UserDataService) ChangeName(ctx context.Context, username string, newName string) (*models.UserInfo, error) {
	const op = "services.user.ChangeName"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("changing name")

	info, err := u.userServiceProvider.ChangeName(ctx, username, newName)
	if err != nil {
		log.Error("failed with changing name", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("name changed")

	return info, nil
}


func (u *UserDataService) ChangePassword(ctx context.Context, username string, oldPassword string, newPassword string) (bool, error) {
	const op = "services.user.ChangePassword"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("changing password")

	oldPassHash, err := u.userServiceProvider.GetPassword(ctx, username)
	if err != nil {
		log.Error("failed to get old password")

		return false, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(oldPassHash, []byte(oldPassword))
	if err != nil {
		log.Error("failed to compare passwords", sl.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	newPasshHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	flag, err := u.userServiceProvider.ChangePassword(ctx, username, newPasshHash)
	if err != nil {
		log.Error("failed to change password")

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("password changed")

	return flag, nil
}

package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	MyProfile(ctx context.Context, username string) (info *models.UserInfo, err error)
	Profile(ctx context.Context, username string) (info *models.UserInfo, err error)
	Password(ctx context.Context, username string) (passHash []byte, err error)
	ListProfiles(ctx context.Context, username string, cursor int64, limit int) (profiles []models.UserInfo, cursors *models.Cursor, err error)
	ChangeUsername(ctx context.Context, oldUsername string, newUsername string) (info *models.UserInfo, err error)
	ChangeEmail(ctx context.Context, username string, newEmail string) (info *models.UserInfo, err error)
	ChangeName(ctx context.Context, username string, newName string) (info *models.UserInfo, err error)
	ChangePassword(ctx context.Context, username string, newPasswordHash []byte) (err error)
}

type UserDataService struct {
	userService UserService
	log         *slog.Logger
}

func New(userProvider UserService, log *slog.Logger) *UserDataService {
	return &UserDataService{
		userService: userProvider,
		log:         log,
	}
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUsersNotFound = errors.New("users not found")
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrPasswordsMismatch = errors.New("passwords don't match")
)

func (u *UserDataService) MyProfile(ctx context.Context, username string) (*models.UserInfo, error) {
	const op = "services.user.MyProfile"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("getting user information")

	info, err := u.userService.MyProfile(ctx, username)
	if err != nil {
		log.Error("failed to get user", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got my profile")

	return info, nil
}

func (u *UserDataService) Profile(ctx context.Context, username string) (*models.UserInfo, error) {
	const op = "services.user.Profile"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("getting profile information")

	info, err := u.userService.Profile(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", sl.Err(err))

			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		log.Error("failed to get profile", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got profile information")

	return info, nil
}

func (u *UserDataService) ListProfiles(ctx context.Context, username string, cursor int64, limit int) ([]models.UserInfo, *models.Cursor, error) {
	const op = "services.user.ListProfiles"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("getting profiles")

	profiles, rcursor, err := u.userService.ListProfiles(ctx, username, cursor, limit)
	if err != nil {
		if errors.Is(err, storage.ErrUsersNotFound) {
			log.Error("users not found")

			return nil, nil, fmt.Errorf("%s: %w", op, ErrUsersNotFound)
		}
		log.Error("failed to get profiles", sl.Err(err))

		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got profiles")

	return profiles, rcursor, nil
}

func (u *UserDataService) UpdateInfo(ctx context.Context, username string, newInfo models.NewUserInfo) (info *models.UserInfo, accessToken string, refreshToken string, err error) {
	const op = "services.user.UpdateInfo"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	if newInfo.NewPassword != "" {
		log.Info("changing password")

		oldPassHash, err := u.userService.Password(ctx, username)
		if err != nil {
			log.Error("failed to get old password")

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		err = bcrypt.CompareHashAndPassword(oldPassHash, []byte(newInfo.PreviousPassword))
		if err != nil {
			log.Error("failed to compare passwords", sl.Err(ErrPasswordsMismatch))

			return nil, "", "", fmt.Errorf("%s: %w", op, ErrPasswordsMismatch)
		}

		newPasshHash, err := bcrypt.GenerateFromPassword([]byte(newInfo.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		err = u.userService.ChangePassword(ctx, username, newPasshHash)
		if err != nil {
			log.Error("failed to change password")

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("password changed")

	}

	if newInfo.Name != "" {
		log.Info("changing name")

		info, err = u.userService.ChangeName(ctx, username, newInfo.Name)
		if err != nil {
			log.Error("failed with changing name", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("name changed")
	}

	if newInfo.Email != "" {
		log.Info("changing email")

		info, err = u.userService.ChangeEmail(ctx, username, newInfo.Email)
		if err != nil {
			if errors.Is(err, storage.ErrEmailExists) {
				log.Error("email already exists", sl.Err(err))

				return nil, "", "", fmt.Errorf("%s: %w", op, ErrEmailExists)
			}
			log.Error("failed to change email", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("email changed")
	}

	if newInfo.Username != "" {
		log.Info("changing username")

		info, err := u.userService.ChangeUsername(ctx, username, newInfo.Username)
		if err != nil {
			if errors.Is(err, storage.ErrUsernameExists) {
				log.Error("username already exists", sl.Err(err))

				return nil, "", "", fmt.Errorf("%s: %w", op, ErrUsernameExists)
			}
			log.Error("failed to change username", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		user := models.InfoToNormalized(info)
		accessToken, refreshToken, err = jwt.NewPairTokens(user)
		if err != nil {
			log.Error("failed to generate tokens", sl.Err(err))

			return nil, "", "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("username changed")
	}

	return info, accessToken, refreshToken, nil
}

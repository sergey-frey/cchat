package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/http-server/handlers"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/cookie"
)

type UserHandlerProvider interface {
	GetUser(ctx context.Context, username string) (info *models.UserInfo, err error)
	ChangeUsername(ctx context.Context, oldUsername string, newUsername string) (info *models.NormalizedUser, accessToken string, refreshToken string, err error)
	ChangeName(ctx context.Context, username string, newName string) (info *models.UserInfo, err error)
	ChangePassword(ctx context.Context, username string, oldPassword string, newPassword string) (flag bool, err error)
}

type UserDataHandler struct {
	userHandlerProvider UserHandlerProvider
	log          *slog.Logger
}

func New(userProvider UserHandlerProvider, log *slog.Logger) *UserDataHandler {
	return &UserDataHandler{
		userHandlerProvider: userProvider,
		log:          log,
	}
}

type newUsername struct {Username string `json:"username"`}
type newName struct {Name string `json:"name"`}
type newPassword struct {OldPassword string `json:"old_password"`; NewPassword string `json:"new_password"`}


func (u *UserDataHandler) GetUser(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.GetUser"

		log := u.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		username, err := cookie.TakeUserInfo(w, r)
		if flag := HandleGettingCookie(w, r, err, log); !flag {
			return
		}

		info, err := u.userHandlerProvider.GetUser(ctx, username)
		if err != nil {
			log.Error("failed to get info")

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to get info",
			})

			return
		}

		log.Info("got info")

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   info,
		})
	}
}

func (u *UserDataHandler) ChangeUsername(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.ChangeUsername"

		log := u.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var newUsername newUsername

		err := render.Decode(r, &newUsername)
		if flag := handlers.HandleError(w, r, newUsername, err, log); !flag {
			return
		}

		oldUsername, err := cookie.TakeUserInfo(w, r)
		if flag := HandleGettingCookie(w, r, err, log); !flag {
			return
		}

		info, accessToken, refreshToken, err := u.userHandlerProvider.ChangeUsername(ctx, oldUsername, newUsername.Username)
		if err != nil {
			log.Error("failed to change username")

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to change username",
			})

			return
		}

		cookie.SetCookie(w, accessToken, refreshToken)

		log.Info("password changed")

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   info,
		})
	}
}

func (u *UserDataHandler) ChangeName(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.ChangeName"

		log := u.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var newName newName

		err := render.Decode(r, &newName)
		if flag := handlers.HandleError(w, r, newName, err, log); !flag {
			return
		}

		username, err := cookie.TakeUserInfo(w, r)
		if flag := HandleGettingCookie(w, r, err, log); !flag {
			return
		}

		info, err := u.userHandlerProvider.ChangeName(ctx, username, newName.Name)
		if err != nil {
			log.Error("failed to change name")

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to change name",
			})

			return
		}

		log.Info("name changed")

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   info,
		})
	}
}

func (u *UserDataHandler) ChangePassword(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.ChangePassword"

		log := u.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var newPassword newPassword

		err := render.Decode(r, &newPassword)
		if flag := handlers.HandleError(w, r, newPassword, err, log); !flag {
			return
		}

		username, err := cookie.TakeUserInfo(w, r)
		if flag := HandleGettingCookie(w, r, err, log); !flag {
			return
		}

		info, err := u.userHandlerProvider.ChangePassword(ctx, username, newPassword.OldPassword, newPassword.NewPassword)
		if err != nil {
			log.Error("failed to change password")

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to change password",
			})

			return
		}

		log.Info("password changed")

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   info,
		})
	}
}

func HandleGettingCookie(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) bool {
	if err != nil {
		log.Error("failed to taking user info")

		render.JSON(w, r, resp.Response{
			Status: http.StatusBadRequest,
			Error:  "failed with getting cookie",
		})

		return false
	}

	return true
}

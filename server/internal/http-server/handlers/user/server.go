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
	UpdateUserInfo(ctx context.Context, username string, newInfo models.NewUserInfo) (info *models.UserInfo, accessToken string, refreshToken string, err error)
}

type UserDataHandler struct {
	userHandlerProvider UserHandlerProvider
	log                 *slog.Logger
}

func New(userProvider UserHandlerProvider, log *slog.Logger) *UserDataHandler {
	return &UserDataHandler{
		userHandlerProvider: userProvider,
		log:                 log,
	}
}


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


func (u *UserDataHandler) UpdateUserInfo(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.UpdateUserInfo"

		log := u.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var newInfo models.NewUserInfo

		err := render.Decode(r, &newInfo)
		if flag := handlers.HandleError(w, r, newInfo, err, log); !flag {
			return
		}

		username, err := cookie.TakeUserInfo(w, r)
		if flag := HandleGettingCookie(w, r, err, log); !flag {
			return
		}

		info, accessToken, refreshToken, err := u.userHandlerProvider.UpdateUserInfo(ctx, username, newInfo)

		if refreshToken != "" {
			cookie.SetCookie(w, accessToken, refreshToken)
		}

		if err != nil {
			log.Error("failed to update user information")

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to update user information",
			})

			return
		}

		log.Info("information changed")

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
package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/http-server/handlers"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/cookie"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/internal/services/user"
)

type User interface {
	GetUser(ctx context.Context, username string) (info *models.UserInfo, err error)
	GetProfile(ctx context.Context, username string) (info *models.UserInfo, err error)
	UpdateUserInfo(ctx context.Context, username string, newInfo models.NewUserInfo) (info *models.UserInfo, accessToken string, refreshToken string, err error)
}

type UserHandler struct {
	userHandler User
	log                 *slog.Logger
}

func New(userProvider User, log *slog.Logger) *UserHandler {
	return &UserHandler{
		userHandler: userProvider,
		log:                 log,
	}
}


// @Summary GetMyProfile
// @Tags user
// @Description Retrieves data about an authenticated user
// @ID get-my-profile
// @Accept  json
// @Produce  json
// @Success 200 {object} models.UserInfo
// @Failure 400,409 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Security CookieAuth
// @Router /user/myprofile [get]
//go:generate go run github.com/vektra/mockery/v2@v2.53 --name=User
func (u *UserHandler) GetUser(ctx context.Context) http.HandlerFunc {
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

		info, err := u.userHandler.GetUser(ctx, username)
		if err != nil {
			log.Error("failed to get info")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
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


// @Summary GetProfile
// @Tags user
// @Description Returns user data, if it exists.
// @ID get-profile
// @Accept  json
// @Produce  json
// @Param username path string true "Existing username"
// @Success 200 {object} models.UserInfo
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Security CookieAuth
// @Router /user/profile/{username} [get]
func (u *UserHandler) GetProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.GetProfile"

		log := u.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		username := chi.URLParam(r, "username")
		if username == "" {
			log.Warn("username is empty")

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.Response{
				Status: http.StatusConflict,
				Error: "invalid request",
			})

			return
		}

		userInfo, err := u.userHandler.GetProfile(ctx, username)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				log.Warn("user not found", "user", username)

				render.Status(r, http.StatusNotFound)

				render.JSON(w, r, resp.Response{
					Status: http.StatusNotFound,
					Error: "user not found",
				})

				return
			}
			log.Error("failed to get profile", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error: "failed to get profile",
			})

			return
		}

		// path := "/profile/" + username

		log.Info("got profile")

		// http.Redirect(w, r, path, http.StatusMovedPermanently)

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: userInfo,
		})
	}
}


// @Summary UpdateProfile
// @Tags user
// @Description Updates the user's information
// @ID update-profile
// @Accept  json
// @Produce  json
// @Param input body models.NewUserInfo true "The new password is at least 8 characters long and has a valid email address."
// @Success 200 {object} models.UserInfo
// @Failure 400,409 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Security CookieAuth
// @Router /user/update [patch]
func (u *UserHandler) UpdateUserInfo(ctx context.Context) http.HandlerFunc {
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

		info, accessToken, refreshToken, err := u.userHandler.UpdateUserInfo(ctx, username, newInfo)

		if refreshToken != "" {
			cookie.SetCookie(w, accessToken, refreshToken)
		}

		if err != nil {
			if errors.Is(err, user.ErrUsernameExists) {
				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "username already exists",
				})

				return
			}

			if errors.Is(err, user.ErrEmailExists) {
				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "email already exists",
				})

				return
			}

			if errors.Is(err, user.ErrPasswordsMismatch) {
				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "passwords don't match",
				})

				return
			}
			
			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to update user information",
			})

			return
		}

		log.Info("information changed successfully")

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   info,
		})
	}
}


func HandleGettingCookie(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) bool {
	if err != nil {
		log.Error("failed to take user info")

		render.Status(r, http.StatusUnauthorized)

		render.JSON(w, r, resp.Response{
			Status: http.StatusUnauthorized,
			Error:  "failed with getting cookie",
		})

		return false
	}

	return true
}
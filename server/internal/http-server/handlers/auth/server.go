package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/http-server/handlers"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/cookie"
	"github.com/sergey-frey/cchat/internal/services/auth"
)

type Auth interface {
	Login(ctx context.Context, loginUser models.LoginUser) (user *models.NormalizedUser, accessToken string, refreshToken string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (user *models.NormalizedUser, accessToken string, refreshToken string, err error)
}

type AuthHandler struct {
	auth Auth
	log  *slog.Logger
}

func New(auth Auth, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		auth: auth,
		log:  log,
	}
}

// @Summary Login
// @Tags auth
// @Description login
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.LoginUser true "account info"
// @Success 200 {object} models.NormalizedUser
// @Failure 400,404,409 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /cchat/auth/login [post]

func (a *AuthHandler) Login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.Login"

		log := a.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.LoginUser

		err := render.Decode(r, &req)
		
		if flag := handlers.HandleError(w, r, req, err, log); !flag {
			return
		}

		user, accessToken, refreshToken, err := a.auth.Login(ctx, req)

		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "invalid email or password",
				})

				return
			}

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})
			return
		}

		cookie.SetCookie(w, accessToken, refreshToken)

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

// @Summary Register
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body todo.RegisterUser true "account info"
// @Success 200 {object} models.NormalizedUser
// @Failure 400,404,409 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /cchat/auth/register [post]

func (a *AuthHandler) Register(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.Register"

		log := a.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.RegisterUser

		err := render.Decode(r, &req)

		if flag := handlers.HandleError(w, r, req, err, log); !flag {
			return
		}

		user, accessToken, refreshToken, err := a.auth.RegisterNewUser(ctx, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrUserExists) {

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "user already exists",
				})

				return
			}

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})
			return
		}

		cookie.SetCookie(w, accessToken, refreshToken)

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

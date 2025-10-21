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
// @Description Accepts email and password and verifies them
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.LoginUser true "valid email and password"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,404,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/login [post]
//
//go:generate go run github.com/vektra/mockery/v2@v2.53 --name=Auth
func (a *AuthHandler) Login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.Login"

		log := a.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.LoginUser

		err := render.DecodeJSON(r.Body, &req)

		if flag := handlers.HandleError(w, r, req, err, log); !flag {
			return
		}

		user, accessToken, refreshToken, err := a.auth.Login(ctx, req)

		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusConflict,
					Error:  "invalid email or password",
				})

				return
			}

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})
			return
		}

		cookie.SetCookie(w, accessToken, refreshToken)

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

// @Summary Register
// @Tags auth
// @Description Accepts the email and password and if the email does not exist creates a new user
// @ID login-account
// @Accept  json
// @Produce  json
// @Param input body models.RegisterUser true "valid email and password(minimum of 8 characters)"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,404,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/register [post]
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

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusConflict,
					Error:  "user already exists",
				})

				return
			}

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})
			return
		}

		cookie.SetCookie(w, accessToken, refreshToken)

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

package auth

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sergey-frey/cchat/internal/domain/models"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/cookie"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/internal/services/auth"
)

type Auth interface {
	Login(ctx context.Context, loginUser models.LoginUser) (user models.NormalizedUser, accessToken string, refreshToken string, err error)
	RegisterNewUser(ctx context.Context, username string, email string, password string) (user models.NormalizedUser, err error)
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
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "empty request",
				})

				return
			}

			log.Error("failed to decode request")

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request",
			})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.ValidationError(validateErr))

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

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error:  "empty request",
				})

				return
			}

			log.Error("failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request",
			})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		user, err := a.auth.RegisterNewUser(ctx, req.Username, req.Email, req.Password)
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

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

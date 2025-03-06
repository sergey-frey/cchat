package session

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/cookie"
	"github.com/sergey-frey/cchat/internal/lib/jwt"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
)

// @Summary Session
// @Tags auth
// @Description check session
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]

func CheckSession(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.CheckSession"

		log = log.With(
			slog.String("op", op),
		)

		flag, err := cookie.CheckCookie(w, r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || errors.Is(err, jwt.ErrUserUnauthorized) {
				log.Warn("user unauthorized", sl.Err(err))

				render.JSON(w, r, resp.Response{
					Status: http.StatusUnauthorized,
					Data: flag,
				})

				return
			}
			
			log.Error("error with check session", sl.Err(err))
			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Data: flag,
	
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: flag,
		})
	}
}

// @Summary Logout
// @Tags auth
// @Description finish session
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]

func FinishSession(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.FinishSession"

		log = log.With(
			"op", op,
		)

		cookie.DeleteCookie(w)

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
		})

	}
}
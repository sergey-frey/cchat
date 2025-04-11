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
// @ID check-session
// @Accept  json
// @Produce  json
// @Success 200 {object} models.NormalizedUser
// @Failure 400,401,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router cchat/auth/session [post]

func CheckSession(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.CheckSession"

		log = log.With(
			slog.String("op", op),
		)

		user, err := cookie.CheckCookie(w, r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || errors.Is(err, jwt.ErrUserUnauthorized) {
				log.Warn("user unauthorized", sl.Err(err))

				render.Status(r, http.StatusUnauthorized)

				render.JSON(w, r, resp.Response{
					Status: http.StatusUnauthorized,
					Error:  "user unauthorized",
				})
				
				return
			}
			
			log.Error("error with check session", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error:  "user unauthorized",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

// @Summary Logout
// @Tags auth
// @Description finish session
// @ID finish-session
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Failure default {object} Response
// @Router cchat/auth/logout[post]

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

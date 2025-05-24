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
// @Description Checks whether a cookie with a token is set
// @ID check-session
// @Accept  json
// @Produce  json
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/session [post]
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

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusUnauthorized,
					Error:  "user unauthorized",
				})
				
				return
			}
			
			log.Error("error with check session", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "user unauthorized",
			})

			return
		}

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

// @Summary Logout
// @Tags auth
// @Description Terminates the user's session, deletes the cookie with the token
// @ID finish-session
// @Accept  json
// @Produce  json
// @Success 200 {object} response.SuccessResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/logout [post]
func FinishSession(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.FinishSession"

		log = log.With(
			"op", op,
		)

		cookie.DeleteCookie(w)

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: "success logout",
		})

	}
}

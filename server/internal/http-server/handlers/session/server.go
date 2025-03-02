package session

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/cookie"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
)

func CheckSession(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.CheckSession"

		log = log.With(
			slog.String("op", op),
		)

		err := cookie.CheckCookie(r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Error("error with check cookie", sl.Err(err))

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Data: "false",
				})

				return
			}
			
			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Data: "false",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: "true",
		})
	}
}

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
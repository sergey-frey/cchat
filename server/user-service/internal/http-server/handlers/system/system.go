package system

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/render"
	resp "github.com/sergey-frey/cchat/user-service/internal/lib/api/response"
)

func GetSystemID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "system.GetSystemID"

		log = log.With(
			slog.String("op", op),
		)

		systemID := os.Getenv("REPLICA_ID")
		if systemID == "" {
			log.Error("failed to get replica id")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to get id of replica",
			})

			return
		}

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   systemID,
		})
	}
}
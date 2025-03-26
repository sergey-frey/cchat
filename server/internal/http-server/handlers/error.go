package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
)

func HandleError(w http.ResponseWriter, r *http.Request, req any, err error, log *slog.Logger) bool {
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.Response{
				Status: http.StatusConflict,
				Error:  "empty request",
			})

			return false
		}

		log.Error("failed to decode request")

		render.Status(r, http.StatusBadRequest)

		render.JSON(w, r, resp.Response{
			Status: http.StatusBadRequest,
			Error:  "failed to decode request",
		})

		return false
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))

		render.Status(r, http.StatusConflict)

		render.JSON(w, r, resp.ValidationError(validateErr))

		return false
	}

	return true
}
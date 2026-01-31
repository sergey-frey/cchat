package chat

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/http-server/handlers"
	"github.com/sergey-frey/cchat/internal/lib/logger/sl"
	resp "github.com/sergey-frey/cchat/internal/lib/api/response"
)

type Chat interface {
	NewChat(ctx context.Context, users []int64) (chatID int64, err error)
}

type ChatRedis interface {
	AddOnline()
	SetOnline()
	SetOfline()
	UpdateOnline()
}

type ChatHandler struct {
	chatHandler Chat
	chatRedisHandler ChatRedis
	log *slog.Logger
}

func New(chatProvider Chat, chatRedisProvider ChatRedis, log *slog.Logger) *ChatHandler {
	return &ChatHandler{
		chatHandler: chatProvider,
		chatRedisHandler: chatRedisProvider,
		log: log,
	}
}

// @Summary NewChat
// @Tags chat
// @Description Creates a new chat
// @ID create-chat
// @Accept  json
// @Produce  json
// @Param input body models.NewChat true "List of users ID's"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Security CookieAuth
// @Router /chat/new [post]
func (ch *ChatHandler) NewChat(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.chat.NewChat"
		
		log := ch.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.NewChat

		err := render.Decode(r, &req)
		if flag := handlers.HandleError(w, r, req, err, log); !flag {
			return
		}

		if len(req.Users) == 0 {
			log.Error("list of users is empty", sl.Err(err))

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error: "list of users is empty",
			})

			return
		}

		chatID, err := ch.chatHandler.NewChat(ctx, req.Users)
		if err != nil {
			log.Error("failed to create new chat", sl.Err(err))

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error: "failed to create new chat",
			})

			return
		}

		log.Info("new chat created", slog.Int64("chat_id", chatID))

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: chatID,
		})
	}
}

func (cs *ChatHandler) AddOnline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (cs *ChatHandler) SetOnline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func (cs *ChatHandler) SetOfline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func (cs *ChatHandler) UpdateOnline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}


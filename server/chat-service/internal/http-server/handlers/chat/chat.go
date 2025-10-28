package chat

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sergey-frey/cchat/server/chat-service/internal/domain/models"
	"github.com/sergey-frey/cchat/server/chat-service/internal/http-server/handlers"
	resp "github.com/sergey-frey/cchat/server/chat-service/internal/lib/api/response"
	"github.com/sergey-frey/cchat/server/chat-service/internal/lib/cookie"
	"github.com/sergey-frey/cchat/server/chat-service/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/server/chat-service/internal/services/chat"
)

type Chat interface {
	NewChat(ctx context.Context, users []int64) (chatID int64, err error)
	ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) (chats []models.Chat, cursors *models.Cursor, err error)
}

type ChatHandler struct {
	chatHandler      Chat
	log              *slog.Logger
}

func New(chatProvider Chat, log *slog.Logger) *ChatHandler {
	return &ChatHandler{
		chatHandler:      chatProvider,
		log:              log,
	}
}

type ChatResponse struct {
	Chats   []models.Chat `json:"chats"`
	RCursor models.Cursor `json:"cursors"`
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
				Error:  "list of users is empty",
			})

			return
		}

		chatID, err := ch.chatHandler.NewChat(ctx, req.Users)
		if err != nil {
			log.Error("failed to create new chat", sl.Err(err))

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to create new chat",
			})

			return
		}

		log.Info("new chat created", slog.Int64("chat_id", chatID))

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   chatID,
		})
	}
}

func (ch *ChatHandler) ListChats(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.chat.Chat"

		log := ch.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userInfo, err := cookie.TakeUserInfo(w, r)
		if flag := handlers.HandleGettingCookie(w, r, err, log); !flag {
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			log.Warn("username is empty")

			render.Status(r, http.StatusOK)

			render.JSON(w, r, resp.SuccessResponse{
				Status: http.StatusOK,
				Data: ChatResponse{
					Chats:   []models.Chat{},
					RCursor: models.Cursor{},
				},
			})

			return
		}

		stringCursor := r.URL.Query().Get("cursor")

		stringPageSize := r.URL.Query().Get("limit")
		if stringPageSize == "" {
			log.Warn("limit is empty")

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusConflict,
				Error:  "limit is empty",
			})

			return
		}

		var cursor int64

		if stringCursor == "" {
			cursor = 0
		} else {
			cursor, err = strconv.ParseInt(stringCursor, 10, 64)
			if err != nil {
				log.Error("failed to convert cursor")

				render.Status(r, http.StatusInternalServerError)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusInternalServerError,
					Error:  "failed to convert curosr",
				})

				return
			}
		}

		limit, err := strconv.Atoi(stringPageSize)
		if err != nil {
			log.Warn("failed to convert limit")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to convert limit",
			})

			return
		}

		if limit < 1 {
			log.Warn("limit must be more than 0")

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "limit must be more than 0",
			})

			return
		}

		chats, cursors, err := ch.chatHandler.ListChats(ctx, userInfo.ID, username, cursor, limit)
		if err != nil {
			if errors.Is(err, chat.ErrChatsNotFound) {
				log.Warn("chats not found")

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusBadRequest,
					Error:  "chats not found",
				})

				return
			}

			log.Error("failed to get chats")

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to get chats",
			})

			return
		}

		log.Info("got chats")

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: ChatResponse{
				Chats:   chats,
				RCursor: *cursors,
			},
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

package message

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sergey-frey/cchat/message-service/internal/domain/models"
	"github.com/sergey-frey/cchat/message-service/internal/http-server/handlers"
	resp "github.com/sergey-frey/cchat/message-service/internal/lib/api/response"
	"github.com/sergey-frey/cchat/message-service/internal/lib/cookie"
	"github.com/sergey-frey/cchat/message-service/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/message-service/internal/services/message"
)

type Message interface {
	NewChat(ctx context.Context, users []int64) (chatID int64, err error)
	ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) (messages []models.Message, cursors *models.Cursor, err error)
}

type MessageHandler struct {
	messageHandler Message
	log            *slog.Logger
}

func New(messageProvider Message, log *slog.Logger) *MessageHandler {
	return &MessageHandler{
		messageHandler: messageProvider,
		log:            log,
	}
}

type HistoryResponse struct {
	ChatID   uuid.UUID        `json:"chat_id"`
	Messages []models.Message `json:"messages"`
	RCursor  models.Cursor    `json:"r_cursor"`
}

func (mh *MessageHandler) SendMessage(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.chat.NewChat"

		log := mh.log.With(
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

		chatID, err := mh.messageHandler.NewChat(ctx, req.Users)
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

func (mh *MessageHandler) ChatHistory(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.message.Chat"

		log := mh.log.With(
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
				Data: HistoryResponse{
					ChatID:   uuid.Nil,
					Messages: []models.Message{},
					RCursor:  models.Cursor{},
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

		messages, cursors, err := mh.messageHandler.ListChats(ctx, userInfo.ID, username, cursor, limit)
		if err != nil {
			if errors.Is(err, message.ErrChatsNotFound) {
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

		//TODO: 
		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: HistoryResponse{
				ChatID: uuid.Nil,
				Messages:   messages,
				RCursor: *cursors,
			},
		})
	}
}

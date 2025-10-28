package chat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sergey-frey/cchat/server/chat-service/internal/domain/models"
	"github.com/sergey-frey/cchat/server/chat-service/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/server/chat-service/internal/provider/storage"
)

type Chat interface {
	NewChat(ctx context.Context, users []int64) (chatID int64, err error)
	ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) (chats []models.Chat, cursors *models.Cursor, err error)
}

type ChatService struct {
	chatService      Chat
	log              *slog.Logger
}

func New(chatProvider Chat, log *slog.Logger) *ChatService {
	return &ChatService{
		chatService:      chatProvider,
		log:              log,
	}
}

var (
	ErrChatsNotFound = errors.New("chats not found")
)

func (cs *ChatService) NewChat(ctx context.Context, users []int64) (chatID int64, err error) {
	const op = "services.chat.NewChat"

	log := cs.log.With(
		slog.String("op", op),
	)

	log.Info("creating chat")

	chatID, err = cs.chatService.NewChat(ctx, users)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return chatID, nil
}

func (cs *ChatService) ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) (chats []models.Chat, cursors *models.Cursor, err error) {
	const op = "services.chat.ListChats"

	log := cs.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("getting chats")

	chats, rcursor, err := cs.chatService.ListChats(ctx, currUser, username, cursor, limit)
	if err != nil {
		if errors.Is(err, storage.ErrChatsNotFound) {
			log.Error("chats not found")

			return nil, nil, fmt.Errorf("%s: %w", op, ErrChatsNotFound)
		}
		log.Error("failed to get chats", sl.Err(err))

		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got chats")

	return chats, rcursor, nil
}

func (cs *ChatService) AddOnline() {

}

func (cs *ChatService) SetOnline() {

}

func (cs *ChatService) SetOfline() {

}

func (cs *ChatService) UpdateOnline() {

}

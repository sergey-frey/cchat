package message

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sergey-frey/cchat/message-service/internal/domain/models"
	"github.com/sergey-frey/cchat/message-service/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/message-service/internal/provider/storage"
)

type Message interface {
	NewChat(ctx context.Context, users []int64) (chatID int64, err error)
	ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) (chats []models.Message, cursors *models.Cursor, err error)
}

type MessageService struct {
	messageService      Message
	log              *slog.Logger
}

func New(messageProvider Message, log *slog.Logger) *MessageService {
	return &MessageService{
		messageService:      messageProvider,
		log:              log,
	}
}

var (
	ErrChatsNotFound = errors.New("chats not found")
)

func (ms *MessageService) NewChat(ctx context.Context, users []int64) (chatID int64, err error) {
	const op = "services.message.NewChat"

	log := ms.log.With(
		slog.String("op", op),
	)

	log.Info("creating chat")

	chatID, err = ms.messageService.NewChat(ctx, users)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return chatID, nil
}

func (ms *MessageService) ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) (chats []models.Message, cursors *models.Cursor, err error) {
	const op = "services.message.ListChats"

	log := ms.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("getting chats")

	chats, rcursor, err := ms.messageService.ListChats(ctx, currUser, username, cursor, limit)
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

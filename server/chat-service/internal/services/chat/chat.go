package chat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sergey-frey/cchat/server/chat-service/internal/domain/models"
	"github.com/sergey-frey/cchat/server/chat-service/internal/lib/logger/sl"
	"github.com/sergey-frey/cchat/server/chat-service/internal/provider/storage"
)

type ChatProvider interface {
	NewChat(ctx context.Context, chatName string,users []uuid.UUID) (chatID uuid.UUID, err error)
	ListChats(ctx context.Context, idUser uuid.UUID, cursor int64, limit int) (chats []models.Chat, cursors *models.Cursor, err error)
}

type UserProvider interface {
	CheckUser(ctx context.Context, ids []uuid.UUID) (error)
}

type ChatService struct {
	chatProvider      ChatProvider
	userProvider	 	UserProvider
	log              *slog.Logger
}

func New(chatProvider ChatProvider, userProvider UserProvider, log *slog.Logger) *ChatService {
	return &ChatService{
		chatProvider:      chatProvider,
		userProvider:    userProvider,
		log:              log,
	}
}

var (
	ErrChatsNotFound = errors.New("chats not found")
)

func (cs *ChatService) NewChat(ctx context.Context, chatName string, users []uuid.UUID) (chatID uuid.UUID, err error) {
	const op = "services.chat.NewChat"

	log := cs.log.With(
		slog.String("op", op),
	)

	log.Info("checking users")

	err = cs.userProvider.CheckUser(ctx, users)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("creating chat")

	chatID, err = cs.chatProvider.NewChat(ctx, chatName,users)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return chatID, nil
}

func (cs *ChatService) ListChats(ctx context.Context, idUser uuid.UUID, cursor int64, limit int) (chats []models.Chat, cursors *models.Cursor, err error) {
	const op = "services.chat.ListChats"

	log := cs.log.With(
		slog.String("op", op),
		slog.String("user_id", idUser.String()),
	)

	log.Info("getting chats")

	chats, rcursor, err := cs.chatProvider.ListChats(ctx, idUser, cursor, limit)
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

package chat

import (
	"context"
	"fmt"
	"log/slog"
)

type Chat interface{
	NewChat(ctx context.Context, users []int64) (chatID int64, err error)
}

type ChatRedis interface {
	AddOnline()
	SetOnline()
	SetOfline()
	UpdateOnline()
}

type ChatService struct {
	chatService Chat
	chatRedisService ChatRedis
	log *slog.Logger
}

func New(chatProvider Chat, chatRedisProvider ChatRedis, log *slog.Logger) *ChatService {
	return &ChatService{
		chatService: chatProvider,
		chatRedisService: chatRedisProvider,
		log: log,
	}
}

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

func (cs *ChatService) AddOnline() {

}

func (cs *ChatService) SetOnline() {

}

func (cs *ChatService) SetOfline() {

}

func (cs *ChatService) UpdateOnline() {
	
}
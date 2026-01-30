package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sergey-frey/cchat/internal/storage"
)

func (s *Storage) NewChat(ctx context.Context, users []int64) (chatID int64, err error) {
	const op = "storage.chat.NewChat"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func ()  {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		
		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	row := s.pool.QueryRow(ctx, `
		INSERT INTO chats DEFAULT VALUES
		RETURNING id;
	`)

	err = row.Scan(&chatID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedToCreateChat)
	}

	rows := make([][]interface{}, len(users))
	for i, userID := range users {
		rows[i] = []interface{}{chatID, userID}
	}

	_, err = s.pool.CopyFrom(
		ctx,
		[]string{"user_chats"},
		[]string{"chat_id", "user_id"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedToAddUsersInChat)
	}

	return chatID, nil
}
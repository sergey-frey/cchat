package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sergey-frey/cchat/message-service/internal/domain/models"
	"github.com/sergey-frey/cchat/message-service/internal/provider/storage"
)

func (s *Storage) NewChat(ctx context.Context, users []int64) (chatID int64, err error) {
	const op = "storage.chat.NewChat"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
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

func (s *Storage) ListChats(ctx context.Context, currUser int64, username string, cursor int64, limit int) ([]models.Message, *models.Cursor, error) {
	const op = "storage.chat.Chat"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	values := make([]interface{}, 0, 5)
	pagination, users, limitQ := "", "", ""
	username += "%"

	if cursor == 0 {
		pagination += fmt.Sprintf("user_id = $%d", len(values)+1)
		users += fmt.Sprintf("$%d", len(values)+2)
		limitQ += fmt.Sprintf("$%d", len(values)+3)
		values = append(values, currUser, username, limit+1)
	}

	if cursor != 0 {
		pagination += fmt.Sprintf("user_id = $%d AND chat_id < $%d", len(values)+1, len(values)+2)
		users += fmt.Sprintf("$%d", len(values)+3)
		limitQ += fmt.Sprintf("$%d", len(values)+4)
		values = append(values, currUser, cursor, username, limit+1)
	}

	stmt := fmt.Sprintf(`
		WITH user_chats_cte AS (
			SELECT chat_id
			FROM user_chats
			WHERE %s
		),
		chat_participants AS (
			SELECT uc.chat_id, u.id AS user_id, u.username, u.name, u.email
			FROM user_chats uc
			JOIN users u ON uc.user_id = u.id
			WHERE uc.chat_id IN (SELECT chat_id FROM user_chats_cte)
		),
		filtered_chats AS (
			SELECT DISTINCT cp.chat_id
			FROM chat_participants cp
			WHERE cp.user_id != $1
			AND (%s::text IS NULL OR cp.username ILIKE %s || '%%')
		),
		last_messages AS (
			SELECT DISTINCT ON (m.chat_id)
				m.chat_id,
				m.id AS message_id,
				m.content,
				m.date,
				m.type,
				m.author_id
			FROM messages m
			WHERE m.chat_id IN (SELECT chat_id FROM filtered_chats)
			ORDER BY m.chat_id, m.date DESC
			LIMIT 1
		)
		SELECT fc.chat_id,
			json_agg(
				json_build_object(
					'id', cp.user_id,
					'username', cp.username,
					'name', cp.name,
					'email', cp.email
				) ORDER BY cp.user_id
			) AS participants,
			json_build_object(
				'id', lm.message_id,
				'content', lm.content,
				'date', lm.date,
				'type', lm.type,
				'author_id', lm.author_id
			) AS last_message
		FROM filtered_chats fc
		JOIN chat_participants cp ON cp.chat_id = fc.chat_id
		LEFT JOIN last_messages lm ON lm.chat_id = fc.chat_id
		GROUP BY fc.chat_id, lm.message_id, lm.content, lm.date, lm.type, lm.author_id
		ORDER BY fc.chat_id DESC
		LIMIT %s;
	`, pagination, users, users, limitQ)

	rows, err := s.pool.Query(ctx, stmt, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("%s: %w", op, storage.ErrChatsNotFound)
		}
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		var userJSON []byte
		var messageJSON []byte

		if err := rows.Scan(&chat.ID, &userJSON, &messageJSON); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}

		if err := json.Unmarshal(userJSON, &chat.Users); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}

		if len(messageJSON) > 0 {
			var msg models.Message
			if err := json.Unmarshal(messageJSON, &msg); err != nil {
				return nil, nil, fmt.Errorf("%s: %w", op, err)
			}
			chat.LastMessage = &msg
		}

		chats = append(chats, chat)
	}

	// rcursor := &models.Cursor{}

	if len(chats) == 0 {
		return nil, nil, fmt.Errorf("%s: %w", op, storage.ErrUsersNotFound)
	}

	return nil, nil, nil
}

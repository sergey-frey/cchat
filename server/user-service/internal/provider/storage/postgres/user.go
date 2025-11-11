package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sergey-frey/cchat/user-service/internal/domain/models"
	"github.com/sergey-frey/cchat/user-service/internal/provider/storage"
)

type pageCursor struct {
	CreatedAt time.Time `json:"c"`
	UUID        string `json:"u"`
}

func (s *Storage) CreateUser(ctx context.Context, uid uuid.UUID, email string, username string, name string) (*models.NormalizedUser, error) {
	const op = "storage.postgres.user.MyProfile"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.NormalizedUser

	fmt.Println(uid, email, username, name)

	row := tx.QueryRow(ctx, `
		INSERT INTO users(user_id, email, username, name)
		VALUES($1, $2, $3, $4)
		RETURNING user_id, email, username;
	`, uid, email, username, name)

	err = row.Scan(&info.UUID, &info.Email, &info.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(info)

	return &info, nil
}

func (s *Storage) GetUserByID(ctx context.Context, username string) (*models.UserInfo, error) {
	const op = "storage.postgres.user.Profile"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.UserInfo

	row := tx.QueryRow(ctx, `
		SELECT id, email, username, name
		FROM users
		WHERE username = $1;
	`, username)

	err = row.Scan(&info.UUID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.NormalizedUser, error) {
	const op = "storage.postgres.user.Profile"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.NormalizedUser

	row := tx.QueryRow(ctx, `
		SELECT user_id, email, username
		FROM users
		WHERE email = $1;
	`, email)

	err = row.Scan(&info.UUID, &info.Email, &info.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) Profiles(ctx context.Context, username string, cursor string, limit int) ([]models.UserInfo, *models.Cursor, error) {
	const op = "storage.postgres.user.ListProfiles"

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

	// Запрашиваем на один элемент больше, чтобы проверить наличие следующей страницы.
	fetchLimit := limit + 1

	var query string
	args := make([]interface{}, 0, 4)

	// Ищем по частичному совпадению имени пользователя.
	// ILIKE - регистронезависимый поиск в PostgreSQL.
	usernamePattern := username + "%"

	// Внутренняя структура для сканирования данных из БД, включая created_at.
	type userWithTimestamp struct {
		models.UserInfo
		CreatedAt time.Time
	}

	if cursor == "" {
		// Первый запрос: курсора нет, начинаем с самого начала.
		// Сортируем по убыванию времени создания, затем по UUID для стабильного порядка.
		query = `
			SELECT id, email, username, name, created_at
			FROM users
			WHERE username ILIKE $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2
		`
		args = append(args, usernamePattern, fetchLimit)
	} else {
		// Последующие запросы: используем курсор.
		decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: invalid cursor format: %w", op, err)
		}

		var c pageCursor
		if err := json.Unmarshal(decodedCursor, &c); err != nil {
			return nil, nil, fmt.Errorf("%s: invalid cursor data: %w", op, err)
		}

		// Ищем записи, которые "старше" (созданы раньше), чем запись в курсоре.
		// Конструкция (created_at, id) < ($2, $3) эффективно использует индекс.
		query = `
			SELECT id, email, username, name, created_at
			FROM users
			WHERE username ILIKE $1 AND (created_at, id) < ($2, $3)
			ORDER BY created_at DESC, id DESC
			LIMIT $4
		`
		args = append(args, usernamePattern, c.CreatedAt, c.UUID, fetchLimit)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		// pgx.ErrNoRows не возвращается для Query, поэтому эта проверка не нужна.
		// Пустой результат - это не ошибка.
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	profiles := make([]userWithTimestamp, 0, fetchLimit)
	for rows.Next() {
		var item userWithTimestamp
		// Сканируем id в поле UUID структуры UserInfo
		err := rows.Scan(&item.UUID, &item.Email, &item.Username, &item.Name, &item.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}
		profiles = append(profiles, item)
	}

	if len(profiles) == 0 {
		return []models.UserInfo{}, &models.Cursor{HasNextPage: false}, nil
	}

	hasNextPage := len(profiles) > limit
	if hasNextPage {
		// Убираем лишний элемент, который мы запрашивали для проверки.
		profiles = profiles[:limit]
	}

	// Создаем курсор для следующей страницы из последнего элемента в текущем списке.
	lastItem := profiles[len(profiles)-1]
	nextCursorData, err := json.Marshal(pageCursor{
		CreatedAt: lastItem.CreatedAt,
		UUID:      lastItem.UUID.String(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%s: failed to create next cursor: %w", op, err)
	}
	
	nextCursor := base64.StdEncoding.EncodeToString(nextCursorData)

	// Преобразуем результат в целевой тип, убирая временную метку.
	resultProfiles := make([]models.UserInfo, len(profiles))
	for i, p := range profiles {
		resultProfiles[i] = p.UserInfo
	}
	
	rcursor := &models.Cursor{
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}

	return resultProfiles, rcursor, nil
}

func (s *Storage) ChangeUsername(ctx context.Context, oldUsername string, newUsername string) (*models.UserInfo, error) {
	const op = "storage.postgres.user.ChangeUsername"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.UserInfo

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET username = $1
		WHERE username = $2
		RETURNING id, email, username, name;
	`, newUsername, oldUsername)

	err = row.Scan(&info.UUID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUsernameExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) ChangeEmail(ctx context.Context, username string, newEmail string) (*models.UserInfo, error) {
	const op = "storage.postgres.user.ChangeEmail"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.UserInfo

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET email = $1
		WHERE username = $2
		RETURNING id, email, username, name;
	`, newEmail, username)

	err = row.Scan(&info.UUID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEmailExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) ChangeName(ctx context.Context, username string, newName string) (*models.UserInfo, error) {
	const op = "storage.postgres.user.ChangeName"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.UserInfo

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET name = $1
		WHERE username = $2
		RETURNING id, email, username, name;
	`, newName, username)

	err = row.Scan(&info.UUID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

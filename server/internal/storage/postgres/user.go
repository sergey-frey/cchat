package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/storage"
)

func (s *Storage) MyProfile(ctx context.Context, username string) (*models.UserInfo, error) {
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

	var info models.UserInfo

	row := tx.QueryRow(ctx, `
		SELECT id, email, username, name
		FROM users
		WHERE username = $1;
	`, username)

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) Profile(ctx context.Context, username string) (*models.UserInfo, error) {
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

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) ListProfiles(ctx context.Context, username string, cursor int64, limit int) ([]models.UserInfo, *models.Cursor, error) {
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

	values := make([]interface{}, 0, 4)
	pagination, limitQ := "", ""
	username += "%"

	if cursor == 0 {
		pagination += fmt.Sprintf("WHERE username LIKE $%d ORDER BY id ASC", len(values)+1)
		limitQ += fmt.Sprintf("$%d", len(values)+2)
		values = append(values, username, limit+1)
	}

	if cursor != 0 {
		pagination += fmt.Sprintf("WHERE id < $%d AND username LIKE $%d ORDER BY id ASC", len(values)+1, len(values)+2)
		limitQ += fmt.Sprintf("$%d", len(values)+3)
		values = append(values, cursor, username, limit+1)
	}

	stmt := fmt.Sprintf(`
		WITH u AS (
			SELECT * FROM users u %s
		)
		SELECT id, email, username, name
		FROM u
		ORDER BY id DESC
		LIMIT %s;
	`, pagination, limitQ)

	rows, err := s.pool.Query(ctx, stmt, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("%s: %w", op, storage.ErrUsersNotFound)
		}
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var profiles []models.UserInfo
	for rows.Next() {
		var item models.UserInfo
		err = rows.Scan(&item.ID, &item.Email, &item.Username, &item.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}
		profiles = append(profiles, item)
	}

	rcursor := &models.Cursor{}

	if len(profiles) == 0 {
		return nil, nil, fmt.Errorf("%s: %w", op, storage.ErrUsersNotFound)
	}

	if len(profiles) < 2 {
		rcursor = &models.Cursor{
			PrevCursor: profiles[len(profiles)-1].ID,
		}
		return profiles, rcursor, nil
	}

	if len(profiles) >= 2 {
		if len(profiles) <= limit {
			rcursor = &models.Cursor{
				PrevCursor: profiles[len(profiles)-1].ID,
			}
			return profiles, rcursor, nil
		}
		if len(profiles) > limit {
			rcursor = &models.Cursor{
				PrevCursor: profiles[len(profiles)-2].ID,
				NextCursor: profiles[len(profiles)-1].ID,
			}
		}
	}

	return profiles[:len(profiles)-1], rcursor, nil
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

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
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

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
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

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) ChangePassword(ctx context.Context, username string, newPasswordHash []byte) error {
	const op = "storage.postgres.user.ChangePassword"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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

	var id int64

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET pass_hash = $1
		WHERE username = $2
		RETURNING id;
	`, newPasswordHash, username)

	err = row.Scan(&id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if id != 0 {
		return nil
	}

	return nil
}

func (s *Storage) Password(ctx context.Context, username string) (password []byte, err error) {
	const op = "storage.postgres.user.Password"

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

	row := tx.QueryRow(ctx, `
		SELECT pass_hash
		FROM users
		WHERE username = $1;
	`, username)

	err = row.Scan(&password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return password, nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/sergey-frey/cchat/internal/domain/models"
	"github.com/sergey-frey/cchat/internal/storage"
)


func (s *Storage) GetUser(ctx context.Context, username string) (*models.UserInfo, error) {
	const op = "storage.postgres.user.GetUser"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

	var info models.UserInfo

	row := tx.QueryRow(ctx, `
		SELECT id, email, username, COALESCE(name, 'nameless') as name
		FROM users
		WHERE username = $1;
	`, username)

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}


func (s *Storage) ChangeUsername(ctx context.Context, oldUsername string,  newUsername string) (*models.UserInfo, error)  {
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
		RETURNING id, email, username, COALESCE(name, 'nameless');
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


func (s *Storage) ChangeName(ctx context.Context, username string, newName string) (*models.UserInfo, error) {
	const op = "storage.postgres.user.ChangeName"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op , err)
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
		RETURNING id, email, username, COALESCE(name, 'nameless');
	`, newName, username)

	err = row.Scan(&info.ID, &info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}


func (s *Storage) ChangePassword(ctx context.Context, username string, newPasswordHash []byte) (error) {
	const op = "storage.postgres.user.ChangePassword"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op , err)
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


func (s *Storage) GetPassword(ctx context.Context, username string) (password []byte, err error) {
	const op = "storage.postgres.user.GetPassword"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op , err)
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
	
	
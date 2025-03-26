package postgres

import (
	"context"
	"fmt"

	"github.com/sergey-frey/cchat/internal/domain/models"
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
		SELECT email, username, COALESCE(name, 'name') as name
		FROM users
		WHERE username = $1;
	`, username)

	err = row.Scan(&info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}


func (s *Storage) ChangeUsername(ctx context.Context, oldUsername string,  newUsername string) (*models.NormalizedUser, error)  {
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

	var user models.NormalizedUser

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET username = $1
		WHERE username = $2
		RETURNING email, username;
	`, newUsername, oldUsername)

	err = row.Scan(&user.Email, &user.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
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
		RETURNING email, username, COALESCE(name, 'nameless');
	`, newName, username)

	err = row.Scan(&info.Email, &info.Username, &info.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}


func (s *Storage) ChangePassword(ctx context.Context, username string, newPasswordHash []byte) (bool, error) {
	const op = "storage.postgres.user.ChangePassword"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op , err)
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
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if id != 0 {
		return true, nil
	}
	
	return false, nil
}


func (s *Storage) GetPassword(ctx context.Context, username string) (password []byte, err error) {
	const op = "storage.postgres.user.GetPassword"

	row := s.pool.QueryRow(ctx, `
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
	
	
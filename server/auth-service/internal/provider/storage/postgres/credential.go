package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) Register(ctx context.Context, uid uuid.UUID, passHash []byte) (error) {
	const op = "storage.postgres.credential.Register"

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

	var rid uuid.UUID

	row := tx.QueryRow(ctx, `
		INSERT INTO credentials(user_id, password_hash, created_at, updated_at)
		VALUES($1, $2, NOW(), NOW())
		RETURNING user_id;
	`, uid, passHash)

	err = row.Scan(&rid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Login(ctx context.Context, uid uuid.UUID) ([]byte, error) {
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

	var password []byte

	row := tx.QueryRow(ctx, `
		SELECT password_hash
		FROM credentials
		WHERE user_id = $1;
	`, uid)

	err = row.Scan(&password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return password, nil
}

func (s *Storage) Password(ctx context.Context, uid uuid.UUID) (password []byte, err error) {
	const op = "storage.postgres.credential.Password"

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
		SELECT password_hash
		FROM credentials
		WHERE user_id = $1;
	`, uid)

	err = row.Scan(&password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return password, nil
}

func (s *Storage) ChangePassword(ctx context.Context, uid uuid.UUID, newPasswordHash []byte) (error) {
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

	var ruid uuid.UUID

	row := tx.QueryRow(ctx, `
		UPDATE credentials
		SET password_hash = $1
		WHERE user_id = $2
		RETURNING user_id;
	`, newPasswordHash, uid)

	err = row.Scan(&ruid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ResetPassword(ctx context.Context, uid uuid.UUID, newPasswordHash []byte) error {
	const op = "storage.postgres.user.ResetPassword"

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

	var ruid uuid.UUID

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET password_hash = $1
		WHERE user_id = $2
		RETURNING user_id;
	`, newPasswordHash, uid)

	err = row.Scan(&ruid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
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

func (s *Storage) SaveUser(ctx context.Context, username string, email string, passHash []byte) ( models.NormalizedUser, error) {
	const op = "storage.postgres.SaveUser"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.NormalizedUser{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	row := tx.QueryRow(ctx, `
		INSERT INTO users (username, email, pass_hash)
		VALUES($1, $2, $3)
		RETURNING id, username, email
	`, username, email, passHash)

	var user models.NormalizedUser

	err = row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.NormalizedUser{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return models.NormalizedUser{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, err
}


func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	row := tx.QueryRow(ctx, `
		SELECT id, email, username, pass_hash
		FROM users
		WHERE email = $1
	`, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
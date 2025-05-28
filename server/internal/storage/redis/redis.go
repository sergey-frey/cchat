package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sergey-frey/cchat/internal/config"
)


type RedisStorage struct {
	client *redis.Client
}

func New(ctx context.Context, cfg config.RedisDB) (*RedisStorage, error) {
	const op = "storage.redis.New"

	db := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		Password: os.Getenv("REDIS_DB_PASSWORD"),
		DB: cfg.DB,
		MaxRetries: cfg.MaxRetries,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RedisStorage{
		client: db,
	}, nil
}
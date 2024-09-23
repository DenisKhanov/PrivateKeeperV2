package cache

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log/slog"
	"time"

	rd "github.com/redis/go-redis/v9"
)

type Redis struct {
	Client     *rd.Client
	setTimeout time.Duration
}

func NewRedis(redisURL, password string, db, setTimeout int) (*Redis, error) {
	client := rd.NewClient(&rd.Options{
		Addr:     redisURL,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	status := client.Ping(ctx)
	if err := status.Err(); err != nil {
		return nil, fmt.Errorf("failed ping redis with url %s %w", redisURL, err)
	}

	logrus.Info("Successful redis connection", slog.String("redis", redisURL))

	return &Redis{
		Client:     client,
		setTimeout: time.Duration(setTimeout) * time.Second,
	}, nil
}

func (r *Redis) HSetWithTTL(key string, data any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.Client.HSet(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash data for key %s: %w", key, err)
	}

	err = r.Client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set TTL for key %s: %w", key, err)
	}

	return nil
}

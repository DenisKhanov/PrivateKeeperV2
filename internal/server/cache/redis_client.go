package cache

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log/slog"
	"time"

	rd "github.com/redis/go-redis/v9"
)

// Redis represents a Redis client with a timeout setting.
type Redis struct {
	Client     *rd.Client    // Redis client for executing commands
	setTimeout time.Duration // Timeout for setting keys in Redis
}

// NewRedis initializes a new Redis client with the provided connection parameters.
// It returns a pointer to the Redis instance and an error if the connection fails.
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

// HSetWithTTL sets a hash value in Redis with a specified TTL (time-to-live).
// It takes a key, the data to be stored, and the TTL duration.
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

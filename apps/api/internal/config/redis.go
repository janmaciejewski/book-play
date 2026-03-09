package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(cfg *RedisConfig) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return RedisClient, nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

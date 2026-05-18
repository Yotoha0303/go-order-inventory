package redis

import (
	"context"
	"fmt"
	"go-order-inventory/config"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {

	password := os.Getenv("REDIS_PASSWORD")

	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis addr missing")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return client, nil
}

package bizcache_test

import (
	"context"
	"go-order-inventory/config"
	"go-order-inventory/global"
	"go-order-inventory/pkg/redis"
	"sync"
	"testing"
	"time"
)

var testRedisOnce sync.Once
var testRedisInitErr error

func setupTestRedis(t *testing.T) {
	t.Helper()

	testRedisOnce.Do(func() {
		cfg, err := config.LoadConfig("../../config.yml")
		if err != nil {
			testRedisInitErr = err
			return
		}

		client, err := redis.InitRedis(cfg.Redis)
		if err != nil {
			testRedisInitErr = err
			return
		}

		global.Redis = client
	})

	if testRedisInitErr != nil {
		t.Skipf("skip redis cache test, init redis failed: %v", testRedisInitErr)
	}

	cleanTestRedis(t)
}

func cleanTestRedis(t *testing.T) {
	t.Helper()

	if global.Redis == nil {
		t.Fatalf("redis client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	keys, err := global.Redis.Keys(ctx, "product:detail:*").Result()
	if err != nil {
		t.Fatalf("list test redis keys failed: %v", err)
	}
	if len(keys) == 0 {
		return
	}

	if err := global.Redis.Del(ctx, keys...).Err(); err != nil {
		t.Fatalf("clean test redis keys failed: %v", err)
	}

	left, err := global.Redis.Keys(ctx, "product:detail:*").Result()
	if err != nil {
		t.Fatalf("verify test redis cleanup failed: %v", err)
	}
	if len(left) != 0 {
		t.Fatalf("expected test redis keys cleaned, remaining=%d", len(left))
	}

}

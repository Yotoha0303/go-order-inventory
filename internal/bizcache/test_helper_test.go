package bizcache_test

import (
	"go-order-inventory/config"
	"go-order-inventory/global"
	"go-order-inventory/pkg/redis"
	"sync"
	"testing"
)

var testRedisOnce sync.Once
var testRedisInitErr error

func setupTestRedis(t *testing.T) {
	t.Helper()

	testRedisOnce.Do(func() {

		cfg, err := config.LoadConfig("../../config.yml")
		if err != nil {
			testRedisInitErr = nil
			return
		}

		client, err := redis.InitRedis(cfg.Redis)
		if err != nil {
			testRedisInitErr = nil
			return
		}

		global.Redis = client
	})

	if testRedisInitErr != nil {
		t.Skipf("skip redis cache test, init redis failed: %v", testRedisInitErr)
	}
}

package bizcache_test

import (
	"go-order-inventory/config"
	"go-order-inventory/global"
	"go-order-inventory/pkg/redis"
	"os"
	"sync"
	"testing"
)

var testRedisOnce sync.Once
var testRedisInitErr error

func setupTestRedis(t *testing.T) {
	t.Helper()

	if os.Getenv("RUN_REDIS_TEST") != "1" {
		t.Skip("skip redis integration test; set RUN_REDIS_TEST=1 to run")
	}

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
}

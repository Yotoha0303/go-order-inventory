package main

import (
	"fmt"
	"go-order-inventory/config"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/pkg/database"
	"go-order-inventory/pkg/redis"
	"go-order-inventory/router"
	"log"
)

var fatalf = log.Fatalf

func main() {
	if err := run(); err != nil {
		fatalf("start server failed: %v", err)
	}
}

func run() error {

	config.LoadEnv()

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		return fmt.Errorf("load config failed:%v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Product{}, &model.Inventory{}, &model.StockLog{}, &model.Order{}, &model.OrderItem{}); err != nil {
		fatalf("auto migrate failed: %v", err)
	}

	global.DB = db

	redisClient, err := redis.InitRedis(cfg)
	if err != nil {
		fatalf("failed to connect redis: %v", err)
	} else {
		global.Redis = redisClient
		fatalf("redis connected")
	}
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Println("server starting at", addr)

	r := router.SetupRouters()

	err = r.Run(addr)
	if err != nil {
		fatalf("run server is failed: %v", err)
	}
}

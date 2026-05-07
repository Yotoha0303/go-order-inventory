package main

import (
	"fmt"
	"go-order-inventory/config"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/pkg/database"
	"go-order-inventory/router"
	"log"
)

func main() {
	config.LoadEnv()

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Product{}, &model.Inventory{}, &model.StockLog{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	global.DB = db

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("load config failed:%v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Println("server starting at", addr)

	r := router.SetupRouters()

	err = r.Run(addr)
	if err != nil {
		log.Fatalf("run server is failed: %v", err)
	}
}

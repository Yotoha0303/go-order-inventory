package app

import (
	"log/slog"

	"go-order-inventory/config"
	"go-order-inventory/internal/bizcache"
	"go-order-inventory/internal/handler"
	"go-order-inventory/internal/service"
	"go-order-inventory/pkg/database"
	"go-order-inventory/pkg/redis"
	"go-order-inventory/router"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	Config      *config.Config
	DB          *gorm.DB
	RedisClient *goredis.Client
	Router      *gin.Engine
	Logger      *slog.Logger
}

func InitDeps(logger *slog.Logger) (*Deps, error) {
	config.LoadEnv()

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		return nil, err
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.InitRedis(cfg)
	if err != nil {
		return nil, err
	}

	productCache := bizcache.NewProductCache(redisClient)

	productService := service.NewProductService(db, productCache)
	inventoryService := service.NewInventoryService(db)
	stockLogService := service.NewStockLogService(db)
	orderService := service.NewOrderService(db)

	handlers := router.Handlers{
		Product:   handler.NewProductHandler(productService),
		Inventory: handler.NewInventoryHandler(inventoryService),
		StockLog:  handler.NewStockLogHandler(stockLogService),
		Order:     handler.NewOrderHandler(orderService),
	}

	r := router.SetupRouters(logger, cfg.HttpServer.Server.Timeout, handlers)

	return &Deps{
		Config:      cfg,
		DB:          db,
		RedisClient: redisClient,
		Router:      r,
		Logger:      logger,
	}, nil
}

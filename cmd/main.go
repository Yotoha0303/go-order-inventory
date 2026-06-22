package main

import (
	"context"
	"fmt"
	"go-order-inventory/config"
	"go-order-inventory/pkg/database"
	"go-order-inventory/pkg/redis"
	"go-order-inventory/router"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type appServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type appDeps struct {
	loadEnv         func()
	loadConfig      func(path string) (*config.Config, error)
	initDB          func(cfg *config.Config) (*gorm.DB, error)
	initRedis       func(cfg *config.Config) (*goredis.Client, error)
	setupRouters    func(db *gorm.DB, logger *slog.Logger, timeout time.Duration, redisClient *goredis.Client) *gin.Engine
	newServer       func(addr string, handler http.Handler, cfg config.HttpServerConfig) appServer
	notify          func(c chan<- os.Signal, sig ...os.Signal)
	shutdownTimeout time.Duration
}

func defaultAppDeps() appDeps {
	return appDeps{
		loadEnv:    config.LoadEnv,
		loadConfig: config.LoadConfig,
		initDB:     database.InitDB,
		initRedis:  redis.InitRedis,
		setupRouters: func(db *gorm.DB, logger *slog.Logger, timeout time.Duration, redisClient *goredis.Client) *gin.Engine {
			return router.SetupRouters(db, logger, timeout, redisClient)
		},
		newServer: func(addr string, handler http.Handler, cfg config.HttpServerConfig) appServer {
			return &http.Server{
				Addr:              addr,
				Handler:           handler,
				ReadTimeout:       cfg.ReadTimeOut,
				WriteTimeout:      cfg.WriteTimeout,
				IdleTimeout:       cfg.IdleTimeout,
				ReadHeaderTimeout: cfg.ReadHeaderTimeout,
				MaxHeaderBytes:    cfg.MaxHeaderBytesKib << 10,
			}
		},
		notify:          signal.Notify,
		shutdownTimeout: 10 * time.Second,
	}
}

var (
	fatalf = log.Fatalf
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func main() {
	if err := run(defaultAppDeps()); err != nil {
		fatalf("start server failed: %v", err)
	}
}

func run(deps appDeps) error {

	deps.loadEnv()

	cfg, err := deps.loadConfig("config.yml")
	if err != nil {
		return fmt.Errorf("load config failed:%v", err)
	}

	db, err := deps.initDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	redisClient, err := deps.initRedis(cfg)
	if err != nil {
		fatalf("failed to connect redis: %v", err)
	} else {
		logger.Info("redis connected")
	}
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	r := deps.setupRouters(db, logger, cfg.HttpServer.Server.Timeout, redisClient)

	server := deps.newServer(addr, r, cfg.HttpServer.Server)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting:", "addr", addr)
		if err := server.ListenAndServe(); err != nil {
			serverErr <- fmt.Errorf("run server failed: %w", err)
			return
		}
		serverErr <- nil
	}()

	quit := make(chan os.Signal, 1)
	deps.notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("server run failed: %w", err)
		}
		return nil
	}

	logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), deps.shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	err = r.Run(addr)
	if err != nil {
		fatalf("run server is failed: %v", err)
	}
	return nil
}

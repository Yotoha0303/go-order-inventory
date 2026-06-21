package main

import (
	"context"
	"fmt"
	"go-order-inventory/config"
	"go-order-inventory/global"
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
)

type appServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

var (
	fatalf = log.Fatalf
)

const (
	shutdownTimeout = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		fatalf("start server failed: %v", err)
	}
}

func newServer(addr string, handler http.Handler, cfg config.HttpServerConfig) appServer {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeOut,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytesKib << 10,
	}
}

func run() error {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config.LoadEnv()

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		return fmt.Errorf("load config failed:%v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	global.DB = db

	redisClient, err := redis.InitRedis(cfg)
	if err != nil {
		fatalf("failed to connect redis: %v", err)
	} else {
		global.Redis = redisClient
		logger.Info("redis connected")
	}
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	r := router.SetupRouters(db, logger, cfg.HttpServer.Server.Timeout)

	server := newServer(addr, r, cfg.HttpServer.Server)

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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("server run failed: %w", err)
		}
		return nil
	}

	logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
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

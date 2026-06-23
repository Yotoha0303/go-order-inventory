package database

import (
	"context"
	"fmt"
	"go-order-inventory/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var openTestMySQL = func(cfg *config.Config, dsn string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql failed: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get mysql db failed: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(cfg.MySQL.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.MySQL.PingTimeout,
	)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return gormDB, nil
}

func buildTestDSN(cfg *config.Config, dbPassword string) (string, error) {
	dbTestDatabse := os.Getenv("MYSQL_TEST_DATABASE")
	if cfg.MySQL.User == "" || dbPassword == "" || cfg.MySQL.Host == "" || cfg.MySQL.Port == "" || dbTestDatabse == "" {
		return "", fmt.Errorf("database config missing")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		dbPassword,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		dbTestDatabse,
	), nil
}

func InitTestDB(cfg *config.Config) (*gorm.DB, error) {
	dsn, err := buildTestDSN(cfg, os.Getenv("MYSQL_TEST_PASSWORD"))
	if err != nil {
		return nil, err
	}

	if dsn == "" {
		return nil, fmt.Errorf("database dsn not found")
	}
	return openTestMySQL(cfg, dsn)
}

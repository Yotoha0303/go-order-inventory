package database

import (
	"fmt"
	"go-order-inventory/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func buildDSN(cfg *config.Config, dbPassword string) (string, error) {
	if cfg.MySQL.User == "" || dbPassword == "" || cfg.MySQL.Host == "" || cfg.MySQL.Port == "" || cfg.MySQL.Database == "" {
		return "", fmt.Errorf("database config missing")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		dbPassword,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	), nil
}

var openMySQL = func(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), nil)
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn, err := buildDSN(cfg, os.Getenv("MYSQL_PASSWORD"))
	if err != nil {
		return nil, err
	}
	return openMySQL(dsn)
}

package database

import (
	"fmt"
	"go-order-inventory/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg config.MySQLConfig) (*gorm.DB, error) {

	dbPassword := os.Getenv("MYSQL_PASSWORD")

	if cfg.User == "" || dbPassword == "" || cfg.Host == "" || cfg.Port == "" || cfg.Database == "" {
		return nil, fmt.Errorf("database config missing")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, dbPassword, cfg.Host, cfg.Port, cfg.Database)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

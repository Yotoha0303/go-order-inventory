package database

import (
	"fmt"
	"go-order-inventory/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTestDB(cfg config.MySQLConfig) (*gorm.DB, error) {

	dbPassword := os.Getenv("MYSQL_TEST_PASSWORD")
	dbDatabase := os.Getenv("MYSQL_TEST_DATABASE")

	if cfg.User == "" || dbPassword == "" || cfg.Host == "" || cfg.Port == "" || dbDatabase == "" {
		return nil, fmt.Errorf("database config missing")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, dbPassword, cfg.Host, cfg.Port, dbDatabase)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTestDB() (*gorm.DB, error) {

	dbUser := os.Getenv("DB_TEST_USER")
	dbPassword := os.Getenv("DB_TEST_PASSWORD")
	dbUrl := os.Getenv("DB_TEST_URL")
	dbPort := os.Getenv("DB_TEST_PORT")
	dbName := os.Getenv("DB_TEST_NAME")

	if dbUser == "" || dbPassword == "" || dbUrl == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("database config missing")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbUrl, dbPort, dbName)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

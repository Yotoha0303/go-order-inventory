package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateStockLog(db *gorm.DB, log *model.StockLog) error {
	return db.Create(log).Error
}

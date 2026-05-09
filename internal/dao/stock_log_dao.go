package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateStockLog(db *gorm.DB, log *model.StockLog) error {
	return db.Create(log).Error
}

func ListStockLogsByProductID(db *gorm.DB, productID int64) ([]*model.StockLog, error) {
	var logs []*model.StockLog
	return logs, db.Where("product_id = ?", productID).Order("created_at desc").First(&logs).Error
}

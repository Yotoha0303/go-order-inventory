package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func InitInventory(db *gorm.DB, inventory *model.Inventory) error {
	return db.Create(inventory).Error
}

func GetInventoryByProductID(db *gorm.DB, productID int64) (*model.Inventory, error) {
	var inventory model.Inventory
	return &inventory, db.Where("product_id = ?", productID).First(&inventory).Error
}

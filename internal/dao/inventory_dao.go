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

func UpdateInventory(db *gorm.DB, productID int64, stockQuantity int64) error {
	var inventory *model.Inventory
	return db.Model(&inventory).Where("product_id = ?", productID).Update("stock_quantity", stockQuantity).Error
}

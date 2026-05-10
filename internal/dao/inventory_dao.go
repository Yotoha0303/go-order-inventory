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

func UpdateInventoryStockQuantity(db *gorm.DB, productID int64, stockQuantity int64) error {

	result := db.Model(&model.Inventory{}).Where("product_id = ?", productID).Update("stock_quantity", stockQuantity)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeductInventory(db *gorm.DB, productID int64, quantity int64) (int64, error) {

	result := db.Model(&model.Inventory{}).Where("product_id = ? AND stock_quantity >= ?", productID, quantity).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", quantity))

	return result.RowsAffected, result.Error
}

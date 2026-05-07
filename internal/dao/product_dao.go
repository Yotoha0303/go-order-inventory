package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB, product *model.Product) error {
	if err := db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

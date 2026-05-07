package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB, product *model.Product) error {
	return db.Create(product).Error
}

package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB, product *model.Product) error {
	return db.Create(product).Error
}

func ListProducts(status int8, db *gorm.DB) ([]*model.Product, error) {
	var products []*model.Product
	return products, db.Model(&model.Product{}).Where("status = ?", status).Order("id DESC").Find(&products).Error
}

func GetProductByID(db *gorm.DB, id int64) (*model.Product, error) {
	var product model.Product
	return &product, db.Where("id = ?", id).First(&product).Error
}

func UpdateProductStatus(db *gorm.DB, id int64, status int8) error {
	return db.Model(&model.Product{}).Where("id = ?", id).Update("status", status).Error
}

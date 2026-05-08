package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB, product *model.Product) error {
	return db.Create(product).Error
}

func ListProducts(status *int8, db *gorm.DB) ([]*model.Product, error) {
	var products []*model.Product
	// if status != 0 {
	// 	return products, db.Find(&products).Error
	// }
	// return products, db.Where("status = ?", status).Find(&products).Error

	query := db.Model(&model.Product{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("id desc").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, err
}

func GetProductByID(db *gorm.DB, id int64) (*model.Product, error) {
	var product model.Product
	return &product, db.Where("id = ?", id).First(&product).Error
}

func UpdateProductStatus(db *gorm.DB, id int64, status int8) error {
	return db.Model(&model.Product{}).Where("id = ?", id).Update("status", status).Error
}

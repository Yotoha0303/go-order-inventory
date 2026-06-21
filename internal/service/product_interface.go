package service

import (
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

type ProductService struct {
	db       *gorm.DB
	daoStore userStore
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db:       db,
		daoStore: daoProductService{},
	}
}

type userStore interface {
	CreateProduct(db *gorm.DB, product *model.Product) error
	ListProducts(db *gorm.DB, status int8) ([]*model.Product, error)
	GetProductByID(db *gorm.DB, id int64) (*model.Product, error)
	UpdateProductStatus(db *gorm.DB, id int64, status int8) error
}

type daoProductService struct{}

func (daoProductService) CreateProduct(db *gorm.DB, product *model.Product) error {
	return dao.CreateProduct(db, product)
}

func (daoProductService) ListProducts(db *gorm.DB, status int8) ([]*model.Product, error) {
	return dao.ListProducts(db, status)
}

func (daoProductService) GetProductByID(db *gorm.DB, id int64) (*model.Product, error) {
	return dao.GetProductByID(db, id)
}

func (daoProductService) UpdateProductStatus(db *gorm.DB, id int64, status int8) error {
	return dao.UpdateProductStatus(db, id, status)
}

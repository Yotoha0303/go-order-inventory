package service

import (
	"context"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

type ProductService struct {
	db       *gorm.DB
	daoStore productStore
	cache    ProductCache
}

func NewProductService(db *gorm.DB, cache ProductCache) *ProductService {
	return &ProductService{
		db:       db,
		daoStore: daoProductService{},
		cache:    cache,
	}
}

type productStore interface {
	CreateProduct(db *gorm.DB, product *model.Product) error
	ListProducts(db *gorm.DB, status int8) ([]*model.Product, error)
	GetProductByID(db *gorm.DB, id int64) (*model.Product, error)
	UpdateProductStatus(db *gorm.DB, id int64, status int8) error
}

type ProductCache interface {
	GetProductDetail(ctx context.Context, productID int64) (*model.Product, bool)
	SetProductDetail(ctx context.Context, product *model.Product)
	DeleteProductDetailCache(ctx context.Context, productID int64)
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

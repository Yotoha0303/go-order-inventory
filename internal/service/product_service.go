package service

import (
	"context"
	"errors"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"strings"

	"gorm.io/gorm"
)

func (p *ProductService) CreateProduct(req request.CreateProductRequest) (*model.Product, error) {
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if req.PriceFen <= 0 {
		return nil, ErrInvalidProductPrice
	}

	if name == "" {
		return nil, ErrInvalidProductName
	}

	if len(description) > 500 {
		return nil, ErrInvalidProductDescription
	}

	product := &model.Product{
		Name:        name,
		Description: description,
		PriceFen:    req.PriceFen,
		Status:      model.ProductStatusOffSale,
	}

	if err := p.daoStore.CreateProduct(p.db, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductService) ListProducts() ([]*model.Product, error) {
	return p.daoStore.ListProducts(p.db, model.ProductStatusOffSale)
}

func (p *ProductService) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {

	if product, ok := p.cache.GetProductDetail(ctx, id); ok {
		return product, nil
	}

	product, err := p.daoStore.GetProductByID(p.db, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	p.cache.SetProductDetail(ctx, product)

	return product, nil
}

func (p *ProductService) OnSaleProduct(ctx context.Context, id int64) error {

	product, err := p.daoStore.GetProductByID(p.db, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err

	}

	if product.Status == model.ProductStatusOnSale {
		return nil
	}

	if err := p.daoStore.UpdateProductStatus(p.db, product.ID, model.ProductStatusOnSale); err != nil {
		return ErrProductOnSaleFailed
	}

	p.cache.DeleteProductDetailCache(ctx, id)

	return nil
}

func (p *ProductService) OffSaleProduct(ctx context.Context, id int64) error {

	product, err := p.daoStore.GetProductByID(p.db, id)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err

	}

	if product.Status == model.ProductStatusOffSale {
		return nil
	}

	if err := p.daoStore.UpdateProductStatus(p.db, product.ID, model.ProductStatusOffSale); err != nil {
		return ErrProductOffSaleFailed
	}

	p.cache.DeleteProductDetailCache(ctx, id)

	return nil
}

package service

import (
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
)

func CreateProduct(req request.CreateProductRequest) (*model.Product, error) {

	return nil, nil
}

func ListProducts(status *int8) ([]*model.Product, error) {

	return nil, nil
}

func GetProductByID(id int64) (*model.Product, error) {

	return nil, nil
}

func OnSaleProduct(id int64) error {

	return nil
}

func OffSaleProduct(id int64) error {

	return nil
}

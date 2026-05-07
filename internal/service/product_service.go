package service

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"strings"
)

var (
	ErrInvalidProductPrice       = errors.New("价格必须大于0")
	ErrInvalidProductName        = errors.New("名称不能为空")
	ErrInvalidProductDescription = errors.New("描述不能超过500个字符")
)

func CreateProduct(req request.CreateProductRequest) (*model.Product, error) {
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

	if err := dao.CreateProduct(global.DB, product); err != nil {
		return nil, err
	}

	return product, nil
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

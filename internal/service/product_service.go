package service

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
)

func CreateProduct(req request.CreateProductRequest) error {

	if req.PriceFen <= 0 {
		return errors.New("价格必须大于0")
	}

	if len(req.Name) == 0 {
		return errors.New("名称不能为空")
	}

	if len(req.Description) > 500 {
		return errors.New("描述不能超过500个字符")
	}

	err := dao.CreateProduct(global.DB, &model.Product{
		Name:        req.Name,
		Description: req.Description,
		PriceFen:    req.PriceFen,
		Status:      2, // 默认下架
	})

	if err != nil {
		return err
	}

	return nil
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

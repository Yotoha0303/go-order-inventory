package service

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrInvalidProductPrice       = errors.New("价格必须大于0")
	ErrInvalidProductName        = errors.New("名称不能为空")
	ErrInvalidProductDescription = errors.New("描述不能超过500个字符")
	ErrProductNotFound           = errors.New("商品信息不存在")
	ErrInvalidProductID          = errors.New("无效的商品ID")
	ErrProductOnSaleFailed       = errors.New("上架商品失败")
	ErrProductAlreadyOnSale      = errors.New("商品已上架")
	ErrProductOffSaleFailed      = errors.New("下架商品失败")
	ErrProductAlreadyOffSale     = errors.New("商品已下架")
	ErrInvalidProductStatus      = errors.New("无效的商品状态")
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
	products, err := dao.ListProducts(status, global.DB)
	if err != nil {
		return nil, err
	}
	return products, err
}

func GetProductByID(id int64) (*model.Product, error) {

	product, err := dao.GetProductByID(global.DB, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return product, nil
}

func OnSaleProduct(id int64) error {

	product, err := dao.GetProductByID(global.DB, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err

	}

	if err := dao.UpdateProductStatus(global.DB, product.ID, model.ProductStatusOnSale); err != nil {
		return ErrProductOnSaleFailed
	}

	return nil
}

func OffSaleProduct(id int64) error {

	product, err := dao.GetProductByID(global.DB, id)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err

	}

	if err := dao.UpdateProductStatus(global.DB, product.ID, model.ProductStatusOffSale); err != nil {
		return ErrProductOffSaleFailed
	}

	return nil
}

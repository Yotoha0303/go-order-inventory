package service

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"

	"gorm.io/gorm"
)

var (
	ErrInitInventoryFailed = errors.New("初始化库存失败")
	ErrInitInventoryExists = errors.New("库存已初始化")
)

func InitInventory(req *request.InitInventoryRequest) error {
	product, err := GetProductByID(req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	data, err := dao.GetInventoryByProductID(global.DB, req.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if data.ID != 0 {
		return ErrInitInventoryExists
	}

	inventory := &model.Inventory{
		ProductID:     product.ID,
		StockQuantity: req.StockQuantity,
	}

	err = dao.InitInventory(global.DB, inventory)
	if err != nil {
		return ErrInitInventoryFailed
	}

	log := &model.StockLog{
		ProductID:      product.ID,
		BeforeQuantity: 0,
		AfterQuantity:  req.StockQuantity,
		ChangeQuantity: req.StockQuantity - 0,
		BizType:        model.StockBizInit,
		Remark:         "初始化库存:" + string(product.Name),
	}

	err = dao.CreateStockLog(global.DB, log)
	if err != nil {
		return ErrCreateStockLogFailed
	}

	return nil
}

func GetInventoryByProductID(productID int64) (*model.Inventory, error) {
	inventory, err := dao.GetInventoryByProductID(global.DB, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return inventory, nil
}

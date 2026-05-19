package service

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/apperror"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrInitInventoryFailed = apperror.New(
		http.StatusInternalServerError,
		response.CodeInitInventoryFailed,
		"初始化库存失败",
	)
	ErrInitInventoryExists = apperror.New(
		http.StatusConflict,
		response.CodeInitInventoryExists,
		"库存已初始化",
	)
	ErrInventoryNotFound = apperror.New(
		http.StatusNotFound,
		response.CodeInventoryNotFound,
		"库存未找到",
	)
	ErrInvalidAddQuantity = apperror.New(
		http.StatusBadRequest,
		response.CodeInventoryInvalidQuantity,
		"增加的库存数量必须大于0",
	)
)

const (
	addInventoryRemarkPrefix  = "手动入库：补充"
	initInventoryRemarkPrefix = "初始化库存："
)

func InitInventory(req *request.InitInventoryRequest) error {
	product, err := dao.GetProductByID(global.DB, req.ProductID)
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

	return global.DB.Transaction(func(tx *gorm.DB) error {
		inventory := &model.Inventory{
			ProductID:     product.ID,
			StockQuantity: *req.StockQuantity,
		}

		if err := dao.InitInventory(tx, inventory); err != nil {
			return ErrInitInventoryFailed
		}

		log := &model.StockLog{
			ProductID:      product.ID,
			BeforeQuantity: 0,
			AfterQuantity:  *req.StockQuantity,
			ChangeQuantity: *req.StockQuantity,
			BizType:        model.StockBizInit,
			Remark:         initInventoryRemarkPrefix + product.Name,
		}

		if err := dao.CreateStockLog(tx, log); err != nil {
			return ErrCreateStockLogFailed
		}
		return nil
	})
}

func GetInventoryByProductID(productID int64) (*model.Inventory, error) {
	inventory, err := dao.GetInventoryByProductID(global.DB, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInventoryNotFound
		}
		return nil, err
	}
	return inventory, nil
}

func AddInventory(req request.AddInventoryRequest) error {
	if req.Quantity <= 0 {
		return ErrInvalidAddQuantity
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {

		inventory, err := dao.GetInventoryByProductIDForUpdate(tx, req.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInventoryNotFound
			}
			return err
		}

		product, err := dao.GetProductByID(tx, req.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProductNotFound
			}
			return err
		}

		beforeQuantity := inventory.StockQuantity
		afterQuantity := beforeQuantity + req.Quantity

		if err := dao.UpdateInventoryStockQuantity(tx, req.ProductID, afterQuantity); err != nil {
			return err
		}

		log := &model.StockLog{
			ProductID:      req.ProductID,
			BeforeQuantity: beforeQuantity,
			AfterQuantity:  afterQuantity,
			ChangeQuantity: req.Quantity,
			BizType:        model.StockBizManualAdd,
			Remark:         addInventoryRemarkPrefix + product.Name,
		}

		err = dao.CreateStockLog(tx, log)
		if err != nil {
			return ErrCreateStockLogFailed
		}
		return nil
	})
}

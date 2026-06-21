package service

import (
	"errors"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"

	"gorm.io/gorm"
)

const (
	addInventoryRemarkPrefix  = "手动入库：补充"
	initInventoryRemarkPrefix = "初始化库存："
)

func (p *InventoryService) InitInventory(req *request.InitInventoryRequest) error {
	if req.StockQuantity == nil {
		return ErrInvalidStockQuantity
	}

	if *req.StockQuantity < 0 {
		return ErrInvalidStockQuantity
	}

	product, err := p.daoStore.GetProductByID(p.db, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	data, err := p.daoStore.GetInventoryByProductID(p.db, req.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if data.ID != 0 {
		return ErrInitInventoryExists
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		inventory := &model.Inventory{
			ProductID:     product.ID,
			StockQuantity: *req.StockQuantity,
		}

		if err := p.daoStore.InitInventory(tx, inventory); err != nil {
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

		if err := p.daoStore.CreateStockLog(tx, log); err != nil {
			return ErrCreateStockLogFailed
		}
		return nil
	})
}

func (p *InventoryService) GetInventoryByProductID(productID int64) (*model.Inventory, error) {
	inventory, err := p.daoStore.GetInventoryByProductID(p.db, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInventoryNotFound
		}
		return nil, err
	}
	return inventory, nil
}

func (p *InventoryService) AddInventory(req request.AddInventoryRequest) error {
	if req.Quantity <= 0 {
		return ErrInvalidAddQuantity
	}

	return p.db.Transaction(func(tx *gorm.DB) error {

		inventory, err := p.daoStore.GetInventoryByProductIDForUpdate(tx, req.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInventoryNotFound
			}
			return err
		}

		product, err := p.daoStore.GetProductByID(tx, req.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProductNotFound
			}
			return err
		}

		beforeQuantity := inventory.StockQuantity
		afterQuantity := beforeQuantity + req.Quantity

		if err := p.daoStore.UpdateInventoryStockQuantity(tx, req.ProductID, afterQuantity); err != nil {
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

		err = p.daoStore.CreateStockLog(tx, log)
		if err != nil {
			return ErrCreateStockLogFailed
		}
		return nil
	})
}

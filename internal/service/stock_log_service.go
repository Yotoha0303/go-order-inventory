package service

import (
	"go-order-inventory/internal/model"
)

func (p *StockLogService) CreateStockLog(log *model.StockLog) error {
	return p.daoStore.CreateStockLog(p.db, log)
}

func (p *StockLogService) ListStockLogsByProductID(productID *int64) ([]*model.StockLog, error) {
	return p.daoStore.ListStockLogsByProductID(p.db, productID)
}

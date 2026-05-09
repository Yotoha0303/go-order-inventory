package service

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
)

var (
	ErrCreateStockLogFailed = errors.New("创建库存日志失败")
	ErrStockLogNotFound     = errors.New("库存日志未找到")
)

func CreateStockLog(log *model.StockLog) error {
	return dao.CreateStockLog(global.DB, log)
}

func ListStockLogsByProductID(productID *int64) ([]*model.StockLog, error) {
	return dao.ListStockLogsByProductID(global.DB, productID)
}

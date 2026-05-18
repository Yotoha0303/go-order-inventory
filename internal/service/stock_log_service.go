package service

import (
	"go-order-inventory/global"
	"go-order-inventory/internal/apperror"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/response"
	"net/http"
)

var (
	ErrCreateStockLogFailed = apperror.New(
		http.StatusNotFound,
		response.CodeCreateStockLogFailed,
		"创建库存日志失败",
	)
)

func CreateStockLog(log *model.StockLog) error {
	return dao.CreateStockLog(global.DB, log)
}

func ListStockLogsByProductID(productID *int64) ([]*model.StockLog, error) {
	return dao.ListStockLogsByProductID(global.DB, productID)
}

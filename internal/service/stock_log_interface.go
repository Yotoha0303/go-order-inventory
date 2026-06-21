package service

import (
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

type StockLogService struct {
	db       *gorm.DB
	daoStore stockLogStore
}

func NewStockLogService(db *gorm.DB) *StockLogService {
	return &StockLogService{
		db:       db,
		daoStore: daoStockLogService{},
	}
}

type stockLogStore interface {
	CreateStockLog(db *gorm.DB, log *model.StockLog) error
	ListStockLogsByProductID(db *gorm.DB, productID *int64) ([]*model.StockLog, error)
}

type daoStockLogService struct{}

func (daoStockLogService) CreateStockLog(db *gorm.DB, log *model.StockLog) error {
	return dao.CreateStockLog(db, log)
}

func (daoStockLogService) ListStockLogsByProductID(db *gorm.DB, productID *int64) ([]*model.StockLog, error) {
	return dao.ListStockLogsByProductID(db, productID)
}

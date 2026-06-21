package service

import (
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

type InventoryService struct {
	db       *gorm.DB
	daoStore inventoryStore
}

func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{
		db:       db,
		daoStore: daoInventoryService{},
	}
}

type inventoryStore interface {
	GetProductByID(db *gorm.DB, id int64) (*model.Product, error)
	InitInventory(db *gorm.DB, inventory *model.Inventory) error
	GetInventoryByProductID(db *gorm.DB, productID int64) (*model.Inventory, error)
	GetInventoryByProductIDForUpdate(db *gorm.DB, productID int64) (*model.Inventory, error)
	UpdateInventoryStockQuantity(db *gorm.DB, productID int64, stockQuantity int64) error
	CreateStockLog(db *gorm.DB, log *model.StockLog) error
}

type daoInventoryService struct{}

func (daoInventoryService) GetProductByID(db *gorm.DB, id int64) (*model.Product, error) {
	return dao.GetProductByID(db, id)
}

func (daoInventoryService) InitInventory(db *gorm.DB, inventory *model.Inventory) error {
	return dao.InitInventory(db, inventory)
}

func (daoInventoryService) GetInventoryByProductID(db *gorm.DB, productID int64) (*model.Inventory, error) {
	return dao.GetInventoryByProductID(db, productID)
}

func (daoInventoryService) GetInventoryByProductIDForUpdate(db *gorm.DB, productID int64) (*model.Inventory, error) {
	return dao.GetInventoryByProductIDForUpdate(db, productID)
}

func (daoInventoryService) UpdateInventoryStockQuantity(db *gorm.DB, productID int64, stockQuantity int64) error {
	return dao.UpdateInventoryStockQuantity(db, productID, stockQuantity)
}

func (daoInventoryService) CreateStockLog(db *gorm.DB, log *model.StockLog) error {
	return dao.CreateStockLog(db, log)
}

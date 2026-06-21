package service

import (
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

type OrderService struct {
	db       *gorm.DB
	daoStore orderStore
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db:       db,
		daoStore: daoOrderService{},
	}
}

type orderStore interface {
	CreateOrder(db *gorm.DB, order *model.Order) error
	CreateOrderItems(db *gorm.DB, items *model.OrderItem) error
	GetOrderByID(db *gorm.DB, id int64) (*model.Order, error)
	ListOrders(db *gorm.DB) ([]*model.Order, error)
	ListOrderItemsByOrderID(db *gorm.DB, orderID int64) ([]*model.OrderItem, error)
	PatchOrderStatus(db *gorm.DB, orderID int64, fromStatus int8, toStatus int8, timeField string) (int64, error)
	PatchOrderTotalPriceFen(db *gorm.DB, orderID int64, totalPriceFen int64) error
	GetProductByID(db *gorm.DB, id int64) (*model.Product, error)
	GetInventoryByProductIDForUpdate(db *gorm.DB, productID int64) (*model.Inventory, error)
	DeductInventory(db *gorm.DB, productID int64, quantity int64) (int64, error)
	UpdateInventoryStockQuantity(db *gorm.DB, productID int64, stockQuantity int64) error
	CreateStockLog(db *gorm.DB, log *model.StockLog) error
}

type daoOrderService struct{}

func (daoOrderService) CreateOrder(db *gorm.DB, order *model.Order) error {
	return dao.CreateOrder(db, order)
}

func (daoOrderService) CreateOrderItems(db *gorm.DB, items *model.OrderItem) error {
	return dao.CreateOrderItems(db, items)
}

func (daoOrderService) GetOrderByID(db *gorm.DB, id int64) (*model.Order, error) {
	return dao.GetOrderByID(db, id)
}

func (daoOrderService) ListOrders(db *gorm.DB) ([]*model.Order, error) {
	return dao.ListOrders(db)
}

func (daoOrderService) ListOrderItemsByOrderID(db *gorm.DB, orderID int64) ([]*model.OrderItem, error) {
	return dao.ListOrderItemsByOrderID(db, orderID)
}

func (daoOrderService) PatchOrderStatus(db *gorm.DB, orderID int64, fromStatus int8, toStatus int8, timeField string) (int64, error) {
	return dao.PatchOrderStatus(db, orderID, fromStatus, toStatus, timeField)
}

func (daoOrderService) PatchOrderTotalPriceFen(db *gorm.DB, orderID int64, totalPriceFen int64) error {
	return dao.PatchOrderTotalPriceFen(db, orderID, totalPriceFen)
}

func (daoOrderService) GetProductByID(db *gorm.DB, id int64) (*model.Product, error) {
	return dao.GetProductByID(db, id)
}

func (daoOrderService) GetInventoryByProductIDForUpdate(db *gorm.DB, productID int64) (*model.Inventory, error) {
	return dao.GetInventoryByProductIDForUpdate(db, productID)
}

func (daoOrderService) DeductInventory(db *gorm.DB, productID int64, quantity int64) (int64, error) {
	return dao.DeductInventory(db, productID, quantity)
}

func (daoOrderService) UpdateInventoryStockQuantity(db *gorm.DB, productID int64, stockQuantity int64) error {
	return dao.UpdateInventoryStockQuantity(db, productID, stockQuantity)
}

func (daoOrderService) CreateStockLog(db *gorm.DB, log *model.StockLog) error {
	return dao.CreateStockLog(db, log)
}

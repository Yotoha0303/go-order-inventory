package dao

import (
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB, order *model.Order) error {
	return db.Create(order).Error
}

func CreateOrderItems(db *gorm.DB, items *model.OrderItem) error {
	return db.Create(items).Error
}

func GetOrderByID(db *gorm.DB, id int64) (*model.Order, error) {
	var order model.Order
	return &order, db.Model(&order).Where("id = ?", id).First(&order).Error
}

func ListOrders(db *gorm.DB) ([]*model.Order, error) {
	var orders []*model.Order
	return orders, db.Model(&model.Order{}).Order("id DESC").Find(&orders).Error
}

func ListOrderItemsByOrderID(db *gorm.DB, orderID int64) ([]*model.OrderItem, error) {
	var items []*model.OrderItem
	return items, db.Model(&model.OrderItem{}).Where("order_id = ?", orderID).Order("id ASC").Find(&items).Error
}

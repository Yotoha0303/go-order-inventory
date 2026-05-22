package service_test

import (
	"go-order-inventory/config"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"go-order-inventory/pkg/database"
	"log"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

var testDBOnce sync.Once
var testDBInitErr error

func setupTestDB(t *testing.T) {
	t.Helper()

	testDBOnce.Do(func() {
		_ = godotenv.Load("../../.env")

		cfg, err := config.LoadConfig("../../config.yml")
		if err != nil {
			log.Fatalf("load config failed:%v", err)
		}

		global.DB, testDBInitErr = database.InitTestDB(cfg.MySQL)
		if testDBInitErr != nil {
			return
		}
		testDBInitErr = global.DB.AutoMigrate(&model.Product{}, &model.Inventory{}, &model.StockLog{}, &model.Order{}, &model.OrderItem{})
	})

	if testDBInitErr != nil {
		t.Skipf("skip integration test, init test db failed: %v", testDBInitErr)
	}

	cleanTables(t)
}

func cleanTables(t *testing.T) {
	t.Helper()

	tables := []string{
		"stock_logs",
		"order_items",
		"orders",
		"product_inventories",
		"products",
	}

	for _, table := range tables {
		if err := global.DB.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatalf("clean table %s failed: %v", table, err)
		}
	}
}

func seedProduct(t *testing.T, name string, priceFen int64, status int8) *model.Product {
	t.Helper()
	p := &model.Product{
		Name:        name,
		Description: name + "-desc",
		PriceFen:    priceFen,
		Status:      status,
	}
	if err := global.DB.Create(p).Error; err != nil {
		t.Fatalf("seed product failed: %v", err)
	}
	return p
}

func seedInventory(t *testing.T, productID int64, qty int64) *model.Inventory {
	t.Helper()
	inv := &model.Inventory{
		ProductID:     productID,
		StockQuantity: qty,
	}
	if err := global.DB.Create(inv).Error; err != nil {
		t.Fatalf("seed inventory failed: %v", err)
	}
	return inv
}

func seedPendingOrder(t *testing.T, name string, status int8, priceFen, qty, orderQty int64) *model.Order {
	t.Helper()

	db := global.DB

	product := &model.Product{
		Name:        name,
		Description: name + "desc",
		Status:      model.ProductStatusOnSale,
		PriceFen:    priceFen,
	}

	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	inventory := &model.Inventory{
		ProductID:     product.ID,
		StockQuantity: qty,
	}

	if err := db.Create(&inventory).Error; err != nil {
		t.Fatalf("create inventory failed: %v", err)
	}

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: product.ID,
				Quantity: orderQty},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	return order
}

func seedPaidOrder(t *testing.T) *model.Order {
	t.Helper()

	db := global.DB
	name := "order paid test"
	priceFen := int64(100)
	qty := int64(50)
	orderQty := int64(25)

	product := &model.Product{
		Name:        name,
		Description: name + "desc",
		Status:      model.ProductStatusOnSale,
		PriceFen:    priceFen,
	}

	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	inventory := &model.Inventory{
		ProductID:     product.ID,
		StockQuantity: qty,
	}

	if err := db.Create(&inventory).Error; err != nil {
		t.Fatalf("create inventory failed: %v", err)
	}

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: product.ID, Quantity: orderQty},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	if err := service.PayOrder(order.ID); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	return order
}

func seedFinishedOrder(t *testing.T) *model.Order {
	t.Helper()

	db := global.DB
	name := "order paid test"
	priceFen := int64(100)
	qty := int64(50)
	orderQty := int64(25)

	product := &model.Product{
		Name:        name,
		Description: name + "desc",
		Status:      model.ProductStatusOnSale,
		PriceFen:    priceFen,
	}

	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	inventory := &model.Inventory{
		ProductID:     product.ID,
		StockQuantity: qty,
	}

	if err := db.Create(&inventory).Error; err != nil {
		t.Fatalf("create inventory failed: %v", err)
	}

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: product.ID, Quantity: orderQty},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	if err := service.PayOrder(order.ID); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	if err := service.FinishOrder(order.ID); err != nil {
		t.Fatalf("finish order failed: %v", err)
	}

	return order
}

func seedCancelledOrder(t *testing.T) *model.Order {
	t.Helper()

	db := global.DB
	name := "order paid test"
	priceFen := int64(100)
	qty := int64(50)
	orderQty := int64(25)

	product := &model.Product{
		Name:        name,
		Description: name + "desc",
		Status:      model.ProductStatusOnSale,
		PriceFen:    priceFen,
	}

	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	inventory := &model.Inventory{
		ProductID:     product.ID,
		StockQuantity: qty,
	}

	if err := db.Create(&inventory).Error; err != nil {
		t.Fatalf("create inventory failed: %v", err)
	}

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: product.ID, Quantity: orderQty},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	if err := service.PayOrder(order.ID); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	if err := service.CancelOrder(order.ID); err != nil {
		t.Fatalf("finish order failed: %v", err)
	}

	return order
}

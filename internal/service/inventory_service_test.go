package service_test

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"testing"
)

func newInventoryService(t *testing.T) *service.InventoryService {
	t.Helper()
	setupTestDB(t)
	return service.NewInventoryService(global.DB)
}

func TestInitInventory_ProductNotFound(t *testing.T) {
	inventorySvc := newInventoryService(t)
	qty := int64(10)
	err := inventorySvc.InitInventory(&request.InitInventoryRequest{
		ProductID:     99999,
		StockQuantity: &qty,
	})
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestInitInventory_Success(t *testing.T) {
	inventorySvc := newInventoryService(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	qty := int64(20)

	err := inventorySvc.InitInventory(&request.InitInventoryRequest{
		ProductID:     p.ID,
		StockQuantity: &qty,
	})
	if err != nil {
		t.Fatalf("init inventory failed: %v", err)
	}

	var inv model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&inv).Error; err != nil {
		t.Fatalf("query inventory failed: %v", err)
	}
	if inv.StockQuantity != qty {
		t.Fatalf("expected stock=%d, got %d", qty, inv.StockQuantity)
	}
}

func TestAddInventory_InvalidQuantity(t *testing.T) {
	inventorySvc := newInventoryService(t)
	err := inventorySvc.AddInventory(request.AddInventoryRequest{
		ProductID: 1,
		Quantity:  0,
	})
	if !errors.Is(err, service.ErrInvalidAddQuantity) {
		t.Fatalf("expected ErrInvalidAddQuantity, got %v", err)
	}
}

func TestAddInventory_Success(t *testing.T) {
	inventorySvc := newInventoryService(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, p.ID, 10)

	err := inventorySvc.AddInventory(request.AddInventoryRequest{
		ProductID: p.ID,
		Quantity:  5,
	})
	if err != nil {
		t.Fatalf("add inventory failed: %v", err)
	}

	var inv model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&inv).Error; err != nil {
		t.Fatalf("query inventory failed: %v", err)
	}
	if inv.StockQuantity != 15 {
		t.Fatalf("expected stock=15, got %d", inv.StockQuantity)
	}
}

func TestInitInventory_CreateStockLog(t *testing.T) {
	inventorySvc := newInventoryService(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	qty := int64(20)

	err := inventorySvc.InitInventory(&request.InitInventoryRequest{
		ProductID:     p.ID,
		StockQuantity: &qty,
	})
	if err != nil {
		t.Fatalf("init inventory failed: %v", err)
	}

	var log model.StockLog
	if err := global.DB.Where("product_id = ?", p.ID).Order("id ASC").First(&log).Error; err != nil {
		t.Fatalf("query stock log failed: %v", err)
	}

	if log.ChangeQuantity != qty || log.AfterQuantity != qty || log.BizType != model.StockBizInit {
		t.Fatalf("unexpected stock log data: %+v", log)
	}
}

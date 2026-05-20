package service_test

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"testing"
)

func TestCreateOrder_InsufficientStock(t *testing.T) {
	setupTestDB(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, p.ID, 1)

	_, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: p.ID, Quantity: 2},
		},
	})
	if !errors.Is(err, service.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}

func TestCreateOrder_Success(t *testing.T) {
	setupTestDB(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, p.ID, 10)

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: p.ID, Quantity: 3},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}
	if order.TotalAmountFen != 300 {
		t.Fatalf("expected total_amount_fen=300, got %d", order.TotalAmountFen)
	}
	if order.Status != model.OrderStatusPending {
		t.Fatalf("unexpected order status: %d", order.Status)
	}

	var inv model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&inv).Error; err != nil {
		t.Fatalf("query inventory failed: %v", err)
	}
	if inv.StockQuantity != 7 {
		t.Fatalf("expected stock=7, got %d", inv.StockQuantity)
	}

	var stockLog model.StockLog
	if err := global.DB.Where("product_id = ? AND biz_id = ?", p.ID, order.ID).First(&stockLog).Error; err != nil {
		t.Fatalf("query stock log failed:%v", err)
	}

	if stockLog.ChangeQuantity != -3 {
		t.Fatalf("expected change_quantity = -3,got %d", stockLog.ChangeQuantity)
	}

	if stockLog.AfterQuantity != 7 {
		t.Fatalf("expected after quantity = 7,got %d", stockLog.AfterQuantity)
	}

	beforeQuantity := stockLog.AfterQuantity + (-stockLog.ChangeQuantity)
	if beforeQuantity != 10 {
		t.Fatalf("expected before quantity = 10, got %d", beforeQuantity)
	}
}

func TestPendingAndPayOrder_Success(t *testing.T) {
	setupTestDB(t)

	var initQuantity int64 = 10
	var initPriceFen int64 = 100

	p := seedProduct(t, "test pending and pay", initPriceFen, model.ProductStatusOnSale)
	seedInventory(t, p.ID, initQuantity)

	var orderedQuantity int64 = 1

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{
				ProductID: p.ID,
				Quantity:  orderedQuantity,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected order created success,%v", err)
	}

	if order.Status != model.OrderStatusPending {
		t.Fatalf("expected status pending,got %d", order.Status)
	}
	var inv model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&inv).Error; err != nil {
		t.Fatalf("expected inventory record exist,got %v", err)
	}

	if inv.StockQuantity != initQuantity-orderedQuantity {
		t.Fatalf("expected stock quantity %d, got %d", initQuantity-orderedQuantity, inv.StockQuantity)
	}

	var stockLogPayBefore model.StockLog
	if err := global.DB.Where("product_id = ? AND biz_id = ?", p.ID, order.ID).First(&stockLogPayBefore).Error; err != nil {
		t.Fatalf("expected inventory record exist,got %v", err)
	}

	if -stockLogPayBefore.ChangeQuantity != orderedQuantity {
		t.Fatalf("expected change quantity equal orderedQuantity,got %d", orderedQuantity)
	}

	if stockLogPayBefore.BizType != model.StockBizOrderDeduct {
		t.Fatalf("expected biz type is 3,got %v", err)
	}

	if err := service.PayOrder(order.ID); err != nil {
		t.Fatalf("expected order pay successes,got %v", err)
	}

	var stockLogPayAfter model.StockLog
	if err := global.DB.Where("product_id = ? AND biz_id = ?", p.ID, order.ID).First(&stockLogPayAfter).Error; err != nil {
		t.Fatalf("expected inventory record exist,got %v", err)
	}

	if stockLogPayAfter.AfterQuantity == stockLogPayBefore.BeforeQuantity {
		t.Fatalf("expected stock log after quantity equal before quantity,got %v", err)
	}

	var orderQuery model.Order
	if err := global.DB.Where("id = ? AND status = ?", order.ID, model.OrderStatusPaid).First(&orderQuery).Error; err != nil {
		t.Fatalf("expected order exist,got %v", err)
	}

}

func TestPayAndFinishOrder_Success(t *testing.T) {
	setupTestDB(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, p.ID, 10)
	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: p.ID, Quantity: 1},
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

	var got model.Order
	if err := global.DB.First(&got, order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}
	if got.Status != model.OrderStatusFinished {
		t.Fatalf("expected finished status, got %d", got.Status)
	}
}

func TestCancelOrder_RollbackInventory(t *testing.T) {
	setupTestDB(t)
	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, p.ID, 10)
	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: p.ID, Quantity: 4},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	if err := service.CancelOrder(order.ID); err != nil {
		t.Fatalf("cancel order failed: %v", err)
	}

	var inv model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&inv).Error; err != nil {
		t.Fatalf("query inventory failed: %v", err)
	}
	if inv.StockQuantity != 10 {
		t.Fatalf("expected stock rollback to 10, got %d", inv.StockQuantity)
	}
}

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

func TestPayOrder_FromPendingToPaid_Success(t *testing.T) {
	setupTestDB(t)

	var initQuantity int64 = 10
	var PriceFen int64 = 100
	var orderedQuantity int64 = 1

	p := seedProduct(t, "pay-order-product", PriceFen, model.ProductStatusOnSale)
	seedInventory(t, p.ID, initQuantity)

	order, err := service.CreateOrder(request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{
				ProductID: p.ID,
				Quantity:  orderedQuantity,
			},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	if order.Status != model.OrderStatusPending {
		t.Fatalf("expected order status pending,got %d", order.Status)
	}
	var invBeforePay model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&invBeforePay).Error; err != nil {
		t.Fatalf("query inventory before pay failed: %v", err)
	}

	expectedStockAfterCreate := initQuantity - orderedQuantity
	if invBeforePay.StockQuantity != expectedStockAfterCreate {
		t.Fatalf("expected stock quantity %d after create order, got %d", expectedStockAfterCreate, invBeforePay.StockQuantity)
	}

	var stockLogCountBeforePay int64
	if err := global.DB.Model(&model.StockLog{}).Where("product_id = ? AND biz_id = ? ", p.ID, order.ID).Count(&stockLogCountBeforePay).Error; err != nil {
		t.Fatalf("count stock logs before pay failed: %v", err)
	}

	if stockLogCountBeforePay != 1 {
		t.Fatalf("expected 1 stock log after create order,got %d", stockLogCountBeforePay)
	}

	if err := service.PayOrder(order.ID); err != nil {
		t.Fatalf("pay order failed:%v", err)
	}

	var paidOrder model.Order
	if err := global.DB.First(&paidOrder, order.ID).Error; err != nil {
		t.Fatalf("query paid order failed:%v", err)
	}

	if paidOrder.Status != model.OrderStatusPaid {
		t.Fatalf("expected order status paid,got %d", paidOrder.Status)
	}

	if paidOrder.PaidAt == nil {
		t.Fatalf("expected paid_at not nil")
	}

	var invAfterPay model.Inventory
	if err := global.DB.Where("product_id = ?", p.ID).First(&invAfterPay).Error; err != nil {
		t.Fatalf("query inventory after pay failed:%v", err)
	}

	if invAfterPay.StockQuantity != expectedStockAfterCreate {
		t.Fatalf("expected stock unchanged after pay:%d,got %d", expectedStockAfterCreate, invAfterPay.StockQuantity)
	}

	var stockLogCountAfterPay int64
	if err := global.DB.Model(&model.StockLog{}).Where("product_id = ? AND biz_id = ?", p.ID, order.ID).Count(&stockLogCountAfterPay).Error; err != nil {
		t.Fatalf("count stock logs after pay failed: %v", err)
	}

	if stockLogCountAfterPay != stockLogCountBeforePay {
		t.Fatalf("expected stock log count unchange after pay,before=%d after=%d", stockLogCountBeforePay, stockLogCountAfterPay)
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

package service_test

import (
	"context"
	"errors"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"sync"
	"testing"

	"gorm.io/gorm"
)

func newOrderService(t *testing.T) (*gorm.DB, *service.OrderService) {
	t.Helper()
	testDB := setupTestDB(t)
	return testDB, service.NewOrderService(testDB)
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	testDB, orderSvc := newOrderService(t)
	p := seedProduct(t, testDB, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, testDB, p.ID, 1)

	_, err := orderSvc.CreateOrder(context.Background(), request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: p.ID, Quantity: 2},
		},
	})
	if !errors.Is(err, service.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}

func TestCreateOrder_Success(t *testing.T) {
	testDB, orderSvc := newOrderService(t)
	p := seedProduct(t, testDB, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, testDB, p.ID, 10)

	order, err := orderSvc.CreateOrder(context.Background(), request.CreateOrderRequest{
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
	if err := testDB.Where("product_id = ?", p.ID).First(&inv).Error; err != nil {
		t.Fatalf("query inventory failed: %v", err)
	}
	if inv.StockQuantity != 7 {
		t.Fatalf("expected stock=7, got %d", inv.StockQuantity)
	}

	var stockLog model.StockLog
	if err := testDB.Where("product_id = ? AND biz_id = ?", p.ID, order.ID).First(&stockLog).Error; err != nil {
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
	testDB, orderSvc := newOrderService(t)

	var initQuantity int64 = 10
	var PriceFen int64 = 100
	var orderedQuantity int64 = 1

	p := seedProduct(t, testDB, "pay-order-product", PriceFen, model.ProductStatusOnSale)
	seedInventory(t, testDB, p.ID, initQuantity)

	order, err := orderSvc.CreateOrder(context.Background(), request.CreateOrderRequest{
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
	if err := testDB.Where("product_id = ?", p.ID).First(&invBeforePay).Error; err != nil {
		t.Fatalf("query inventory before pay failed: %v", err)
	}

	expectedStockAfterCreate := initQuantity - orderedQuantity
	if invBeforePay.StockQuantity != expectedStockAfterCreate {
		t.Fatalf("expected stock quantity %d after create order, got %d", expectedStockAfterCreate, invBeforePay.StockQuantity)
	}

	var stockLogCountBeforePay int64
	if err := testDB.Model(&model.StockLog{}).Where("product_id = ? AND biz_id = ? ", p.ID, order.ID).Count(&stockLogCountBeforePay).Error; err != nil {
		t.Fatalf("count stock logs before pay failed: %v", err)
	}

	if stockLogCountBeforePay != 1 {
		t.Fatalf("expected 1 stock log after create order,got %d", stockLogCountBeforePay)
	}

	if err := orderSvc.PayOrder(context.Background(), order.ID); err != nil {
		t.Fatalf("pay order failed:%v", err)
	}

	var paidOrder model.Order
	if err := testDB.First(&paidOrder, order.ID).Error; err != nil {
		t.Fatalf("query paid order failed:%v", err)
	}

	if paidOrder.Status != model.OrderStatusPaid {
		t.Fatalf("expected order status paid,got %d", paidOrder.Status)
	}

	if paidOrder.PaidAt == nil {
		t.Fatalf("expected paid_at not nil")
	}

	var invAfterPay model.Inventory
	if err := testDB.Where("product_id = ?", p.ID).First(&invAfterPay).Error; err != nil {
		t.Fatalf("query inventory after pay failed:%v", err)
	}

	if invAfterPay.StockQuantity != expectedStockAfterCreate {
		t.Fatalf("expected stock unchanged after pay:%d,got %d", expectedStockAfterCreate, invAfterPay.StockQuantity)
	}

	var stockLogCountAfterPay int64
	if err := testDB.Model(&model.StockLog{}).Where("product_id = ? AND biz_id = ?", p.ID, order.ID).Count(&stockLogCountAfterPay).Error; err != nil {
		t.Fatalf("count stock logs after pay failed: %v", err)
	}

	if stockLogCountAfterPay != stockLogCountBeforePay {
		t.Fatalf("expected stock log count unchange after pay,before=%d after=%d", stockLogCountBeforePay, stockLogCountAfterPay)
	}

}

func TestPayAndFinishOrder_Success(t *testing.T) {
	testDB, orderSvc := newOrderService(t)
	p := seedProduct(t, testDB, "p1", 100, model.ProductStatusOnSale)
	seedInventory(t, testDB, p.ID, 10)
	order, err := orderSvc.CreateOrder(context.Background(), request.CreateOrderRequest{
		Items: []request.CreateOrderItemRequest{
			{ProductID: p.ID, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	if err := orderSvc.PayOrder(context.Background(), order.ID); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}
	if err := orderSvc.FinishOrder(context.Background(), order.ID); err != nil {
		t.Fatalf("finish order failed: %v", err)
	}

	var got model.Order
	if err := testDB.First(&got, order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}
	if got.Status != model.OrderStatusFinished {
		t.Fatalf("expected finished status, got %d", got.Status)
	}
}

func TestPayOrder_AlreadyPaid_ReturnsError(t *testing.T) {
	testDB, orderSvc := newOrderService(t)

	order := seedPaidOrder(t, testDB)

	err := orderSvc.PayOrder(context.Background(), order.ID)
	if !errors.Is(err, service.ErrOrderAlreadyPaid) {
		t.Fatalf("expected ErrOrderAlreadyPaid,got %v", err)
	}

	var got model.Order
	if err := testDB.First(&got, order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}

	if got.Status != model.OrderStatusPaid {
		t.Fatalf("expected order status still paid,got %d", got.Status)
	}
}

func TestFinishOrder_PendingOrder_ReturnsNotPaidError(t *testing.T) {
	testDB, orderSvc := newOrderService(t)

	order := seedPendingOrder(t, testDB)

	err := orderSvc.FinishOrder(context.Background(), order.ID)
	if !errors.Is(err, service.ErrOrderNotPaid) {
		t.Fatalf("expected order unpaid is not finished,got %v", err)
	}

	var got model.Order
	if err := testDB.First(&got, order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}

	if got.Status != model.OrderStatusPending {
		t.Fatalf("expected order status unpaid,got %d", got.Status)
	}

	if got.CompletedAt != nil {
		t.Fatalf("expected completed_at nil,got %v", got.CompletedAt)
	}
}

func TestCancelOrder_Success(t *testing.T) {
	testDB, orderSvc := newOrderService(t)

	db := testDB
	ctx := seedPendingOrderContext(t, testDB)

	if err := orderSvc.CancelOrder(context.Background(), ctx.Order.ID); err != nil {
		t.Fatalf("expected order cancel success,got %v", err)
	}

	var got model.Order
	if err := db.First(&got, ctx.Order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}

	if got.Status != model.OrderStatusCancelled {
		t.Fatalf("expected order status already cancel,got %d", got.Status)
	}

	if got.CancelledAt == nil {
		t.Fatalf("expected order cancelled_at not null,got %v", got.CancelledAt)
	}

	var inv model.Inventory
	if err := db.Where("product_id = ?", ctx.Product.ID).First(&inv).Error; err != nil {
		t.Fatalf("query inventory failed: %v", err)
	}

	if inv.StockQuantity != ctx.InitQty {
		t.Fatalf("expected product inventory already rollback,got %d", inv.StockQuantity)
	}

	var rollbackLog model.StockLog
	if err := db.Where("product_id = ? AND biz_id = ? AND biz_type = ?", ctx.Product.ID, ctx.Order.ID, model.StockBizOrderRollback).Order("created_at DESC").First(&rollbackLog).Error; err != nil {
		t.Fatalf("query stock log failed: %v", err)
	}

	if rollbackLog.ChangeQuantity != ctx.OrderQty {
		t.Fatalf("expected rollback change_quantity=%d,got %d", ctx.OrderQty, rollbackLog.ChangeQuantity)
	}

	if rollbackLog.AfterQuantity != ctx.InitQty {
		t.Fatalf("expected stock log after_quantity=%d,got %d", ctx.InitQty, rollbackLog.AfterQuantity)
	}

}

func TestPaidOrder_UnableCancel_ReturnsError(t *testing.T) {
	testDB, orderSvc := newOrderService(t)

	order := seedPaidOrder(t, testDB)
	db := testDB

	err := orderSvc.CancelOrder(context.Background(), order.ID)
	if !errors.Is(err, service.ErrOrderAlreadyPaid) {
		t.Fatalf("expected order cancel failed,got %v", err)
	}

	var got model.Order
	if err := db.First(&got, order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}

	if got.Status == model.OrderStatusCancelled {
		t.Fatalf("expected order status is %d,got %d", model.OrderStatusPaid, got.Status)
	}

	if got.CancelledAt != nil {
		t.Fatalf("expected order cancel failed,got %v", got.CancelledAt)
	}
}

func TestFinishedOrder_UnableCancel_ReturnsError(t *testing.T) {
	testDB, orderSvc := newOrderService(t)

	db := testDB

	order := seedFinishedOrder(t, testDB)

	err := orderSvc.CancelOrder(context.Background(), order.ID)
	if !errors.Is(err, service.ErrOrderAlreadyFinished) {
		t.Fatalf("expected order cancel failed,got %d", err)
	}

	var got model.Order
	if err := db.First(&got, order.ID).Error; err != nil {
		t.Fatalf("query order failed: %v", err)
	}

	if got.Status == model.OrderStatusCancelled {
		t.Fatalf("expected order status cancelled,got %d", got.Status)
	}

	if got.CancelledAt != nil {
		t.Fatalf("expected order cancel failed,got %v", got.CancelledAt)
	}
}

func TestCancelOrder_CancelledOrder_Idempotent(t *testing.T) {
	testDB, orderSvc := newOrderService(t)

	ctx := seedPendingOrderContext(t, testDB)
	db := testDB

	if err := orderSvc.CancelOrder(context.Background(), ctx.Order.ID); err != nil {
		t.Fatalf("order cancenl failed: %v", err)
	}

	var invAfterFirstCancel model.Inventory
	if err := db.Where("product_id = ?", ctx.Product.ID).First(&invAfterFirstCancel).Error; err != nil {
		t.Fatalf("query inventory after first cancel failed: %v", err)
	}

	var rollbackLogCountAfterFirstCancel int64
	if err := db.Model(&model.StockLog{}).Where("product_id = ? AND biz_id = ? AND biz_type = ?", ctx.Product.ID, ctx.Order.ID, model.StockBizOrderRollback).Count(&rollbackLogCountAfterFirstCancel).Error; err != nil {
		t.Fatalf("count rollback stock logs after first cancel failed: %v", err)
	}

	var orderAfterFirstCancel model.Order
	if err := db.First(&orderAfterFirstCancel, ctx.Order.ID).Error; err != nil {
		t.Fatalf("query order after first cancel failed: %v", err)
	}
	if orderAfterFirstCancel.Status != model.OrderStatusCancelled {
		t.Fatalf("expected order status cancelled after first cancel, got %d", orderAfterFirstCancel.Status)
	}
	if orderAfterFirstCancel.CancelledAt == nil {
		t.Fatalf("expected cancelled_at not nil after first cancel")
	}

	if err := orderSvc.CancelOrder(context.Background(), ctx.Order.ID); err != nil {
		t.Fatalf("order cancenl failed: %v", err)
	}

	var invAfterSecondCancel model.Inventory
	if err := db.Where("product_id = ?", ctx.Product.ID).First(&invAfterSecondCancel).Error; err != nil {
		t.Fatalf("query inventory after second cancel failed: %v", err)
	}

	if invAfterSecondCancel.StockQuantity != invAfterFirstCancel.StockQuantity {
		t.Fatalf("expected inventory unchanged on second cancel, before=%d after=%d", invAfterFirstCancel.StockQuantity, invAfterSecondCancel.StockQuantity)
	}

	var rollbackLogCountAfterSecondCancel int64
	if err := db.Model(&model.StockLog{}).Where("product_id = ? AND biz_id = ? AND biz_type = ?", ctx.Product.ID, ctx.Order.ID, model.StockBizOrderRollback).Count(&rollbackLogCountAfterSecondCancel).Error; err != nil {
		t.Fatalf("count rollback stock logs after second cancel failed: %v", err)
	}

	if rollbackLogCountAfterFirstCancel != rollbackLogCountAfterSecondCancel {
		t.Fatalf("expected rollback log count unchanged on second cancel, before=%d after=%d", rollbackLogCountAfterFirstCancel, rollbackLogCountAfterSecondCancel)
	}

	var orderAfterSecondCancel model.Order
	if err := db.First(&orderAfterSecondCancel, ctx.Order.ID).Error; err != nil {
		t.Fatalf("query order after second cancel failed: %v", err)
	}
	if orderAfterSecondCancel.Status != model.OrderStatusCancelled {
		t.Fatalf("expected order status still cancelled after second cancel, got %d", orderAfterSecondCancel.Status)
	}
	if orderAfterSecondCancel.CancelledAt == nil {
		t.Fatalf("expected cancelled_at not nil after second cancel")
	}

}

func TestOrder_ConcurrentTesting_OrderOversold(t *testing.T) {
	const (
		initialStock = int64(10)
		requests     = 20
	)

	testDB, orderSvc := newOrderService(t)
	product := seedProduct(t, testDB, "test concurrent testing order", 100, model.ProductStatusOnSale)

	seedInventory(t, testDB, product.ID, initialStock)

	var wg sync.WaitGroup
	start := make(chan struct{})
	errCh := make(chan error, requests)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			_, err := orderSvc.CreateOrder(context.Background(), request.CreateOrderRequest{
				Items: []request.CreateOrderItemRequest{
					{
						ProductID: product.ID,
						Quantity:  1,
					},
				},
			})
			errCh <- err
		}()
	}

	close(start)
	wg.Wait()
	close(errCh)

	var successCount, failedCount int64
	for err := range errCh {
		if err == nil {
			successCount++
			continue
		}
		if !errors.Is(err, service.ErrInsufficientStock) {
			t.Fatalf("expected ErrInsufficientStock, got %v", err)
		}
		failedCount++
	}

	if successCount != initialStock {
		t.Fatalf("expected success count %d, got %d", initialStock, successCount)
	}
	if failedCount != requests-initialStock {
		t.Fatalf("expected failed count %d, got %d", requests-initialStock, failedCount)
	}

	var stockLogCount int64
	if err := testDB.Model(&model.StockLog{}).Where("product_id = ? AND biz_type = ?", product.ID, model.StockBizOrderDeduct).Count(&stockLogCount).Error; err != nil {
		t.Fatalf("count stock logs failed: %v", err)
	}

	if stockLogCount != initialStock {
		t.Fatalf("expected stock log count %d,got %d", initialStock, failedCount)
	}

	var inventoryFinaly model.Inventory
	if err := testDB.Where("product_id = ?", product.ID).First(&inventoryFinaly).Error; err != nil {
		t.Fatalf("query inventory failed:%v", err)
	}

	if inventoryFinaly.StockQuantity != 0 {
		t.Fatalf("expected inventory quantity is 0,got %d", inventoryFinaly.StockQuantity)
	}

}

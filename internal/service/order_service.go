package service

import (
	"context"
	"errors"
	"fmt"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"time"

	"gorm.io/gorm"
)

func (p *OrderService) CreateOrder(ctx context.Context, req request.CreateOrderRequest) (*model.Order, error) {
	var createOrder *model.Order

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalAmountFen int64 = 0

		orderNo := generateOrderNo()

		order := &model.Order{
			OrderNo:        orderNo,
			TotalAmountFen: 0,
			Status:         model.OrderStatusPending,
		}

		if err := dao.CreateOrder(ctx, tx, order); err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			product, err := dao.GetProductByID(ctx, tx, itemReq.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrProductNotFound
				}
				return err
			}

			if product.Status != model.ProductStatusOnSale {
				return ErrProductOffSale
			}

			inv, err := dao.GetInventoryByProductIDForUpdate(ctx, tx, itemReq.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrInventoryNotFound
				}
				return err
			}

			beforeQuantity := inv.StockQuantity
			afterQuantity := beforeQuantity - itemReq.Quantity

			if afterQuantity < 0 {
				return ErrInsufficientStock
			}

			rows, err := dao.DeductInventory(ctx, tx, itemReq.ProductID, itemReq.Quantity)
			if err != nil {
				return err
			}
			if rows == 0 {
				return ErrInsufficientStock
			}

			subtotalFen := product.PriceFen * itemReq.Quantity
			totalAmountFen += subtotalFen

			orderItem := &model.OrderItem{
				OrderID:         order.ID,
				ProductID:       product.ID,
				ProductName:     product.Name,
				ProductPriceFen: product.PriceFen,
				Quantity:        itemReq.Quantity,
				SubtotalFen:     subtotalFen,
			}

			if err := dao.CreateOrderItems(ctx, tx, orderItem); err != nil {
				return err
			}

			stockLog := &model.StockLog{
				ProductID:      product.ID,
				ChangeQuantity: -itemReq.Quantity,
				BeforeQuantity: beforeQuantity,
				AfterQuantity:  afterQuantity,
				BizType:        model.StockBizOrderDeduct,
				BizID:          &order.ID,
				Remark:         "创建订单扣减库存：" + order.OrderNo,
			}
			if err := dao.CreateStockLog(ctx, tx, stockLog); err != nil {
				return ErrCreateStockLogFailed
			}
		}

		if err := dao.PatchOrderTotalPriceFen(ctx, tx, order.ID, totalAmountFen); err != nil {
			return err
		}

		order.TotalAmountFen = totalAmountFen
		createOrder = order

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createOrder, nil
}

func generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

func (p *OrderService) GetOrderByID(ctx context.Context, id int64) (*model.Order, []*model.OrderItem, error) {
	order, err := dao.GetOrderByID(ctx, p.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrOrderNotFound
		}
		return nil, nil, err
	}

	items, err := dao.ListOrderItemsByOrderID(ctx, p.db, id)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

func (p *OrderService) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return dao.ListOrders(ctx, p.db)
}

func (p *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		order, err := dao.GetOrderByID(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderNotFound
			}
			return err
		}

		switch order.Status {
		case model.OrderStatusCancelled:
			return nil
		case model.OrderStatusPending:
		case model.OrderStatusPaid:
			return ErrOrderAlreadyPaid
		case model.OrderStatusFinished:
			return ErrOrderAlreadyFinished
		default:
			return ErrOrderCancelFailed
		}

		rows, err := dao.PatchOrderStatus(ctx, tx, order.ID, model.OrderStatusPending, model.OrderStatusCancelled, "cancelled_at")
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrOrderCancelFailed
		}

		items, err := dao.ListOrderItemsByOrderID(ctx, tx, order.ID)
		if err != nil {
			return err
		}

		for _, item := range items {
			inventory, err := dao.GetInventoryByProductIDForUpdate(ctx, tx, item.ProductID)
			if err != nil {
				return err
			}

			before := inventory.StockQuantity
			after := before + item.Quantity

			if err := dao.UpdateInventoryStockQuantity(ctx, tx, item.ProductID, after); err != nil {
				return err
			}

			stockLog := &model.StockLog{
				ProductID:      item.ProductID,
				BizID:          &order.ID,
				ChangeQuantity: item.Quantity,
				AfterQuantity:  after,
				BeforeQuantity: before,
				BizType:        model.StockBizOrderRollback,
				Remark:         "取消订单回滚库存：" + order.OrderNo,
			}

			if err := dao.CreateStockLog(ctx, tx, stockLog); err != nil {
				return err
			}
		}

		return nil
	})
}

func (p *OrderService) PayOrder(ctx context.Context, orderID int64) error {
	order, err := dao.GetOrderByID(ctx, p.db, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return ErrOrderAlreadyPaid
	case model.OrderStatusFinished:
		return ErrOrderAlreadyFinished
	case model.OrderStatusCancelled:
		return ErrOrderAlreadyCanceled
	case model.OrderStatusPending:
	default:
		return ErrOrderPayFailed
	}

	row, err := dao.PatchOrderStatus(ctx, p.db, order.ID, model.OrderStatusPending, model.OrderStatusPaid, "paid_at")
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrOrderPayFailed
	}
	return nil
}

func (p *OrderService) FinishOrder(ctx context.Context, orderID int64) error {
	order, err := dao.GetOrderByID(ctx, p.db, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	switch order.Status {
	case model.OrderStatusPending:
		return ErrOrderNotPaid
	case model.OrderStatusCancelled:
		return ErrOrderAlreadyCanceled
	case model.OrderStatusFinished:
		return ErrOrderAlreadyFinished
	case model.OrderStatusPaid:
	default:
		return ErrOrderFinishFailed
	}

	row, err := dao.PatchOrderStatus(ctx, p.db, order.ID, model.OrderStatusPaid, model.OrderStatusFinished, "completed_at")
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrOrderFinishFailed
	}
	return nil

}

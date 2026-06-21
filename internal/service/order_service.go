package service

import (
	"errors"
	"fmt"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"time"

	"gorm.io/gorm"
)

func (p *OrderService) CreateOrder(req request.CreateOrderRequest) (*model.Order, error) {
	var createOrder *model.Order

	err := p.db.Transaction(func(tx *gorm.DB) error {
		var totalAmountFen int64 = 0

		orderNo := generateOrderNo()

		order := &model.Order{
			OrderNo:        orderNo,
			TotalAmountFen: 0,
			Status:         model.OrderStatusPending,
		}

		if err := p.daoStore.CreateOrder(tx, order); err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			product, err := p.daoStore.GetProductByID(tx, itemReq.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrProductNotFound
				}
				return err
			}

			if product.Status != model.ProductStatusOnSale {
				return ErrProductOffSale
			}

			inv, err := p.daoStore.GetInventoryByProductIDForUpdate(tx, itemReq.ProductID)
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

			rows, err := p.daoStore.DeductInventory(tx, itemReq.ProductID, itemReq.Quantity)
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

			if err := p.daoStore.CreateOrderItems(tx, orderItem); err != nil {
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
			if err := p.daoStore.CreateStockLog(tx, stockLog); err != nil {
				return ErrCreateStockLogFailed
			}
		}

		if err := p.daoStore.PatchOrderTotalPriceFen(tx, order.ID, totalAmountFen); err != nil {
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

func (p *OrderService) GetOrderByID(id int64) (*model.Order, []*model.OrderItem, error) {
	order, err := p.daoStore.GetOrderByID(p.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrOrderNotFound
		}
		return nil, nil, err
	}

	items, err := p.daoStore.ListOrderItemsByOrderID(p.db, id)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

func (p *OrderService) ListOrders() ([]*model.Order, error) {
	return p.daoStore.ListOrders(p.db)
}

func (p *OrderService) CancelOrder(orderID int64) error {
	return p.db.Transaction(func(tx *gorm.DB) error {

		order, err := p.daoStore.GetOrderByID(tx, orderID)
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

		rows, err := p.daoStore.PatchOrderStatus(tx, order.ID, model.OrderStatusPending, model.OrderStatusCancelled, "cancelled_at")
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrOrderCancelFailed
		}

		items, err := p.daoStore.ListOrderItemsByOrderID(tx, order.ID)
		if err != nil {
			return err
		}

		for _, item := range items {
			inventory, err := p.daoStore.GetInventoryByProductIDForUpdate(tx, item.ProductID)
			if err != nil {
				return err
			}

			before := inventory.StockQuantity
			after := before + item.Quantity

			if err := p.daoStore.UpdateInventoryStockQuantity(tx, item.ProductID, after); err != nil {
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

			if err := p.daoStore.CreateStockLog(tx, stockLog); err != nil {
				return err
			}
		}

		return nil
	})
}

func (p *OrderService) PayOrder(orderID int64) error {
	order, err := p.daoStore.GetOrderByID(p.db, orderID)
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

	row, err := p.daoStore.PatchOrderStatus(p.db, order.ID, model.OrderStatusPending, model.OrderStatusPaid, "paid_at")
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrOrderPayFailed
	}
	return nil
}

func (p *OrderService) FinishOrder(orderID int64) error {
	order, err := p.daoStore.GetOrderByID(p.db, orderID)
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

	row, err := p.daoStore.PatchOrderStatus(p.db, order.ID, model.OrderStatusPaid, model.OrderStatusFinished, "completed_at")
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrOrderFinishFailed
	}
	return nil

}

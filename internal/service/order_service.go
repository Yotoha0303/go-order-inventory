package service

import (
	"errors"
	"fmt"
	"go-order-inventory/global"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"time"

	"gorm.io/gorm"
)

var (
	ErrProductOffSale      = errors.New("商品已下架")
	ErrInsufficientStock   = errors.New("库存不足")
	ErrCreateOrderFailed   = errors.New("创建订单失败")
	ErrOrderNotFound       = errors.New("订单不存在")
	ErrOrderCancelFailed   = errors.New("订单取消失败")
	ErrOrderPayFailed      = errors.New("订单支付失败")
	ErrOrderFinishFailed   = errors.New("订单完成失败")
	ErrOrderPendingFailed  = errors.New("订单未支付")
	ErrOrderAlreadCanceled = errors.New("订单已取消")
	ErrOrderAlreadFinished = errors.New("订单已完成")
	ErrOrderAlreadPaid     = errors.New("订单已支付")
)

func CreateOrder(req request.CreateOrderRequest) (*model.Order, error) {
	var createOrder *model.Order

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmountFen int64 = 0

		orderNo := generateOrderNo()

		order := &model.Order{
			OrderNo:        orderNo,
			TotalAmountFen: 0,
			Status:         model.OrderStatusPending,
		}

		if err := dao.CreateOrder(tx, order); err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			product, err := dao.GetProductByID(tx, itemReq.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrProductNotFound
				}
				return err
			}

			if product.Status != model.ProductStatusOnSale {
				return ErrProductOffSale
			}

			inv, err := dao.GetInventoryByProductID(tx, itemReq.ProductID)
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

			rows, err := dao.DeductInventory(tx, itemReq.ProductID, itemReq.Quantity)
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

			if err := dao.CreateOrderItems(tx, orderItem); err != nil {
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
			if err := dao.CreateStockLog(tx, stockLog); err != nil {
				return ErrCreateStockLogFailed
			}
		}

		if err := tx.Model(&model.Order{}).Where("id = ?", order.ID).Update("total_amount_fen", totalAmountFen).Error; err != nil {
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

func GetOrderByID(id int64) (*response.OrderDetailResponse, error) {
	order, err := dao.GetOrderByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	items, err := dao.ListOrderItemsByOrderID(global.DB, id)
	if err != nil {
		return nil, err
	}

	return &response.OrderDetailResponse{
		Order: order,
		Items: items,
	}, nil
}

func ListOrders() ([]*model.Order, error) {
	return dao.ListOrders(global.DB)
}

func CancelOrders(orderID int64) error {
	order, err := dao.GetOrderByID(global.DB, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status == model.OrderStatusCancelled {
		return ErrOrderAlreadCanceled
	}

	if order.Status == model.OrderStatusFinished {
		return ErrOrderAlreadFinished
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		row, err := dao.PatchOrderStatus(tx, order.ID, model.OrderStatusPending, model.OrderStatusCancelled, "cancelled_at")
		if err != nil {
			return err
		}

		if row == 0 {
			return ErrOrderCancelFailed
		}

		orderItems, err := dao.ListOrderItemsByOrderID(tx, order.ID)
		if err != nil {
			return err
		}

		stockLogData, err := dao.ListStockLogsByProductID(tx, &orderItems[len(orderItems)-1].ProductID)
		if err != nil {
			return err
		}

		changeQuantity := stockLogData[0].BeforeQuantity - stockLogData[0].AfterQuantity
		inventory, err := dao.GetInventoryByProductID(tx, stockLogData[0].ProductID)
		if err != nil {
			return err
		}
		sum := inventory.StockQuantity + changeQuantity
		err = dao.UpdateInventoryStockQuantity(tx, stockLogData[0].ProductID, sum)
		if err != nil {
			return err
		}

		stockLog := &model.StockLog{
			ProductID:      stockLogData[0].ProductID,
			BizID:          &order.ID,
			ChangeQuantity: changeQuantity,
			AfterQuantity:  sum,
			BeforeQuantity: stockLogData[0].AfterQuantity,
			BizType:        model.StockBizOrderRollback,
			Remark:         "取消订单加库存：" + order.OrderNo,
		}

		err = dao.CreateStockLog(tx, stockLog)
		if err != nil {
			return err
		}
		return nil
	})
}

func PayOrder(orderID int64) error {
	order, err := dao.GetOrderByID(global.DB, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status == model.OrderStatusPaid {
		return ErrOrderAlreadPaid
	}

	if order.Status == model.OrderStatusFinished {
		return ErrOrderAlreadFinished
	}

	if order.Status == model.OrderStatusCancelled {
		return ErrOrderAlreadCanceled
	}

	if order.Status == model.OrderStatusPending {
		row, err := dao.PatchOrderStatus(global.DB, order.ID, model.OrderStatusPending, model.OrderStatusPaid, "paid_at")
		if err != nil {
			return err
		}

		if row == 0 {
			return ErrOrderPayFailed
		}
		return nil
	}
	return ErrOrderPayFailed
}

func FinishOrder(orderID int64) error {
	order, err := dao.GetOrderByID(global.DB, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status == model.OrderStatusPending {
		return ErrOrderPendingFailed
	}

	if order.Status == model.OrderStatusCancelled {
		return ErrOrderAlreadCanceled
	}

	if order.Status == model.OrderStatusFinished {
		return ErrOrderAlreadFinished
	}

	if order.Status == model.OrderStatusPaid {
		row, err := dao.PatchOrderStatus(global.DB, order.ID, model.OrderStatusPaid, model.OrderStatusFinished, "completed_at")
		if err != nil {
			return err
		}

		if row == 0 {
			return ErrOrderFinishFailed
		}
		return nil
	}
	return ErrOrderFinishFailed
}

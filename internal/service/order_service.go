package service

import (
	"errors"
	"fmt"
	"go-order-inventory/global"
	"go-order-inventory/internal/apperror"
	"go-order-inventory/internal/dao"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"net/http"
	"time"

	"gorm.io/gorm"
)

var (
	ErrProductOffSale = apperror.New(
		http.StatusConflict,
		response.CodeProductAlreadyOffSale,
		"商品已下架",
	)

	ErrInsufficientStock = apperror.New(
		http.StatusConflict,
		response.CodeInsufficientStock,
		"库存不足",
	)

	ErrCreateOrderFailed = apperror.New(
		http.StatusInternalServerError,
		response.CodeCreateOrderFailed,
		"创建订单失败",
	)

	ErrOrderNotFound = apperror.New(
		http.StatusNotFound,
		response.CodeOrderNotFound,
		"订单不存在",
	)

	ErrOrderPayFailed = apperror.New(
		http.StatusConflict,
		response.CodeOrderPayFailed,
		"订单支付失败",
	)

	ErrOrderFinishFailed = apperror.New(
		http.StatusConflict,
		response.CodeOrderFinishFailed,
		"订单完成失败",
	)

	ErrOrderCancelFailed = apperror.New(
		http.StatusConflict,
		response.CodeOrderCancelFailed,
		"订单取消失败",
	)

	ErrOrderNotPaid = apperror.New(
		http.StatusConflict,
		response.CodeOrderNotPaid,
		"订单未支付",
	)

	ErrOrderAlreadyCanceled = apperror.New(
		http.StatusConflict,
		response.CodeOrderAlreadyCanceled,
		"订单已取消",
	)

	ErrOrderAlreadyFinished = apperror.New(
		http.StatusConflict,
		response.CodeOrderAlreadyFinished,
		"订单已完成",
	)

	ErrOrderAlreadyPaid = apperror.New(
		http.StatusConflict,
		response.CodeOrderAlreadyPaid,
		"订单已支付",
	)
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

			inv, err := dao.GetInventoryByProductIDForUpdate(tx, itemReq.ProductID)
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

		if err := dao.PatchOrderTotalPriceFen(tx, order.ID, totalAmountFen); err != nil {
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

func GetOrderByID(id int64) (*model.Order, []*model.OrderItem, error) {
	order, err := dao.GetOrderByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrOrderNotFound
		}
		return nil, nil, err
	}

	items, err := dao.ListOrderItemsByOrderID(global.DB, id)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

func ListOrders() ([]*model.Order, error) {
	return dao.ListOrders(global.DB)
}

func CancelOrder(orderID int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {

		order, err := dao.GetOrderByID(tx, orderID)
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

		rows, err := dao.PatchOrderStatus(tx, order.ID, model.OrderStatusPending, model.OrderStatusCancelled, "cancelled_at")
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrOrderCancelFailed
		}

		items, err := dao.ListOrderItemsByOrderID(tx, order.ID)
		if err != nil {
			return err
		}

		for _, item := range items {
			inventory, err := dao.GetInventoryByProductIDForUpdate(tx, item.ProductID)
			if err != nil {
				return err
			}

			before := inventory.StockQuantity
			after := before + item.Quantity

			if err := dao.UpdateInventoryStockQuantity(tx, item.ProductID, after); err != nil {
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

			if err := dao.CreateStockLog(tx, stockLog); err != nil {
				return err
			}
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

	row, err := dao.PatchOrderStatus(global.DB, order.ID, model.OrderStatusPending, model.OrderStatusPaid, "paid_at")
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrOrderPayFailed
	}
	return nil
}

func FinishOrder(orderID int64) error {
	order, err := dao.GetOrderByID(global.DB, orderID)
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

	row, err := dao.PatchOrderStatus(global.DB, order.ID, model.OrderStatusPaid, model.OrderStatusFinished, "completed_at")
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrOrderFinishFailed
	}
	return nil

}

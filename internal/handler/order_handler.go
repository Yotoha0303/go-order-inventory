package handler

import (
	"errors"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 3000, "请求参数错误")
		return
	}

	order, err := service.CreateOrder(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, http.StatusNotFound, 3001, err.Error())
		case errors.Is(err, service.ErrProductOffSale):
			response.Fail(c, http.StatusConflict, 3002, err.Error())
		case errors.Is(err, service.ErrInventoryNotFound):
			response.Fail(c, http.StatusNotFound, 3003, err.Error())
		case errors.Is(err, service.ErrInsufficientStock):
			response.Fail(c, http.StatusConflict, 3004, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 3005, err.Error())
		}
		return
	}

	response.Success(c, order)
}

func GetOrderByID(c *gin.Context) {
	id, ok := parsePositiveProductID(c, "id")
	if !ok {
		return
	}

	order, err := service.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Fail(c, http.StatusNotFound, 3007, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, 3008, "查询订单详情失败")
		return
	}
	response.Success(c, order)
}

func ListOrders(c *gin.Context) {
	orders, err := service.ListOrders()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 3006, "查询订单列表失败")
		return
	}
	response.Success(c, orders)
}

func CancelOrders(c *gin.Context) {
	orderID, ok := parsePositiveProductID(c, "id")
	if !ok {
		return
	}

	if err := service.CancelOrders(orderID); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) || errors.Is(err, service.ErrOrderCancelFailed) || errors.Is(err, service.ErrOrderAlreadPayFinished) {
			response.Fail(c, http.StatusBadRequest, 3009, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, 3010, "取消订单失败")
		return
	}

	response.Success(c, nil)
}

func PayOrder(c *gin.Context) {
	orderID, ok := parsePositiveProductID(c, "id")
	if !ok {
		return
	}

	if err := service.PayOrder(orderID); err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			response.Fail(c, http.StatusBadRequest, 3011, err.Error())
		case errors.Is(err, service.ErrOrderPayFailed), errors.Is(err, service.ErrOrderAlreadCancel), errors.Is(err, service.ErrOrderAlreadFinished):
			response.Fail(c, http.StatusBadRequest, 3012, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 3010, "订单支付失败")
			return
		}
	}

	response.Success(c, nil)
}

func FinishOrder(c *gin.Context) {
	orderID, ok := parsePositiveProductID(c, "id")
	if !ok {
		return
	}

	if err := service.FinishOrder(orderID); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) ||
			errors.Is(err, service.ErrOrderCancelFailed) {
			response.Fail(c, http.StatusBadRequest, 3012, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, 3010, "订单出现错误")
		return
	}

	response.Success(c, nil)
}

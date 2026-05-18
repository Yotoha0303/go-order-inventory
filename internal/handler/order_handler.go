package handler

import (
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err, http.StatusBadRequest, "请求参数错误")
		return
	}

	order, err := service.CreateOrder(req)
	if err != nil {
		handleError(c, err, response.CodeCreateOrderFailed, "订单创建失败")
		return
	}
	response.Success(c, order)
}

func ListOrders(c *gin.Context) {
	orders, err := service.ListOrders()
	if err != nil {
		handleError(c, err, response.CodeQueryOrderListFailed, "查询订单列表失败")
		return
	}
	response.Success(c, orders)
}

func GetOrderByID(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	order, err := service.GetOrderByID(id)
	if err != nil {
		handleError(c, err, response.CodeQueryOrderDetailNotFound, "查询订单详情失败")
		return
	}
	response.Success(c, order)
}

func PayOrder(c *gin.Context) {
	orderID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	if err := service.PayOrder(orderID); err != nil {
		handleError(c, err, response.CodeOrderPayFailed, "支付订单失败")
		return
	}

	response.Success(c, nil)
}

func FinishOrder(c *gin.Context) {
	orderID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	if err := service.FinishOrder(orderID); err != nil {
		handleError(c, err, response.CodeOrderFinishConflict, "完成订单失败")
		return
	}

	response.Success(c, nil)
}

func CancelOrders(c *gin.Context) {
	orderID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	if err := service.CancelOrder(orderID); err != nil {
		handleError(c, err, response.CodeOrderCancelConflict, "取消订单失败")
		return
	}

	response.Success(c, nil)
}

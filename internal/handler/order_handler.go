package handler

import (
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderService interface {
	CreateOrder(req request.CreateOrderRequest) (*model.Order, error)
	ListOrders() ([]*model.Order, error)
	GetOrderByID(id int64) (*model.Order, []*model.OrderItem, error)
	PayOrder(orderID int64) error
	FinishOrder(orderID int64) error
	CancelOrder(orderID int64) error
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

var _ OrderService = (*service.OrderService)(nil)

func (p *OrderHandler) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeOrderParameterError, "请求参数错误")
		return
	}

	order, err := p.orderService.CreateOrder(req)
	if err != nil {
		handleError(c, err, response.CodeCreateOrderFailed, "订单创建失败")
		return
	}

	response.Success(c, order)
}

func (p *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := p.orderService.ListOrders()
	if err != nil {
		handleError(c, err, response.CodeQueryOrderListFailed, "查询订单列表失败")
		return
	}

	response.Success(c, orders)
}

func (p *OrderHandler) GetOrderByID(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	order, orderItems, err := p.orderService.GetOrderByID(id)
	if err != nil {
		handleError(c, err, response.CodeQueryOrderDetailFailed, "查询订单详情失败")
		return
	}

	orderDetail := response.OrderDetailResponse{
		Order: order,
		Items: orderItems,
	}

	response.Success(c, orderDetail)
}

func (p *OrderHandler) PayOrder(c *gin.Context) {
	orderID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	if err := p.orderService.PayOrder(orderID); err != nil {
		handleError(c, err, response.CodeOrderPayFailed, "支付订单失败")
		return
	}

	response.Success(c, nil)
}

func (p *OrderHandler) FinishOrder(c *gin.Context) {
	orderID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	if err := p.orderService.FinishOrder(orderID); err != nil {
		handleError(c, err, response.CodeOrderFinishFailed, "完成订单失败")
		return
	}

	response.Success(c, nil)
}

func (p *OrderHandler) CancelOrders(c *gin.Context) {
	orderID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}

	if err := p.orderService.CancelOrder(orderID); err != nil {
		handleError(c, err, response.CodeOrderCancelFailed, "取消订单失败")
		return
	}

	response.Success(c, nil)
}

package handler

import (
	"context"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req request.CreateOrderRequest) (*model.Order, error)
	ListOrders(ctx context.Context) ([]*model.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*model.Order, []*model.OrderItem, error)
	PayOrder(ctx context.Context, orderID int64) error
	FinishOrder(ctx context.Context, orderID int64) error
	CancelOrder(ctx context.Context, orderID int64) error
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

	if req.IdempotencyKey == "" {
		response.Fail(c, http.StatusBadRequest, response.CodeOrderParameterError, "idempotency key 不能为空")
		return
	}

	order, err := p.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		handleError(c, err, response.CodeCreateOrderFailed, "订单创建失败")
		return
	}

	response.Success(c, order)
}

func (p *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := p.orderService.ListOrders(c.Request.Context())
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

	order, orderItems, err := p.orderService.GetOrderByID(c.Request.Context(), id)
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

	if err := p.orderService.PayOrder(c.Request.Context(), orderID); err != nil {
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

	if err := p.orderService.FinishOrder(c.Request.Context(), orderID); err != nil {
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

	if err := p.orderService.CancelOrder(c.Request.Context(), orderID); err != nil {
		handleError(c, err, response.CodeOrderCancelFailed, "取消订单失败")
		return
	}

	response.Success(c, nil)
}

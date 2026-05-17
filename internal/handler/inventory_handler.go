package handler

import (
	"errors"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitInventory(c *gin.Context) {
	var req request.InitInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParameterError, "请求参数错误")
		return
	}

	if err := service.InitInventory(&req); err != nil {
		switch {
		case err == service.ErrProductNotFound:
			response.Fail(c, http.StatusNotFound, response.CodeProductNotFound, err.Error())
		case err == service.ErrInitInventoryFailed:
			response.Fail(c, http.StatusInternalServerError, response.CodeInitInventoryFailed, err.Error())
		case err == service.ErrInitInventoryExists:
			response.Fail(c, http.StatusConflict, response.CodeInitInventoryExists, err.Error())
		case errors.Is(err, service.ErrCreateStockLogFailed):
			response.Fail(c, http.StatusInternalServerError, response.CodeCreateStockLogFailed, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, response.CodeInitInventoryFailed, "未知错误")
		}
		return
	}

	response.Success(c, nil)
}

func AddInventory(c *gin.Context) {
	var req request.AddInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParameterError, "请求参数错误")
		return
	}
	if err := service.AddInventory(req); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAddQuantity):
			response.Fail(c, http.StatusBadRequest, response.CodeInvalidAddQuantityFailed, err.Error())

		case errors.Is(err, service.ErrInventoryNotFound):
			response.Fail(c, http.StatusNotFound, response.CodeInventoryNotFound, err.Error())

		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, http.StatusNotFound, response.CodeProductNotFound, err.Error())

		case errors.Is(err, service.ErrCreateStockLogFailed):
			response.Fail(c, http.StatusInternalServerError, response.CodeCreateStockLogFailed, err.Error())

		default:
			response.Fail(c, http.StatusInternalServerError, response.CodeAddInventoryError, "增加库存失败")
		}
		return
	}
	response.Success(c, nil)
}

func GetInventoryByProductID(c *gin.Context) {
	id, ok := parsePositiveID(c, "product_id")
	if !ok {
		return
	}

	inventory, err := service.GetInventoryByProductID(id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInventoryNotFound, "查询库存失败")
		return
	}

	response.Success(c, inventory)
}

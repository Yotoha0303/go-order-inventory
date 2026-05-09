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
		response.Fail(c, http.StatusBadRequest, 2000, "请求参数错误")
		return
	}

	if err := service.InitInventory(&req); err != nil {
		switch {
		case err == service.ErrProductNotFound:
			response.Fail(c, http.StatusNotFound, 2001, err.Error())
		case err == service.ErrInitInventoryFailed:
			response.Fail(c, http.StatusInternalServerError, 2002, err.Error())
		case err == service.ErrInitInventoryExists:
			response.Fail(c, http.StatusConflict, 2003, err.Error())
		case errors.Is(err, service.ErrCreateStockLogFailed):
			response.Fail(c, http.StatusInternalServerError, 2004, err.Error())
		case errors.Is(err, service.ErrInventoryNotFound):
			response.Fail(c, http.StatusNotFound, 2005, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 2006, "未知错误")
		}
		return
	}

	response.Success(c, nil)
}

func AddInventory(c *gin.Context) {
	var req request.AddInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 2000, "请求参数错误")
		return
	}
	if err := service.AddInventory(req); err != nil {
		response.Fail(c, http.StatusInternalServerError, 2001, err.Error())
		return
	}
	response.Success(c, nil)
}

func GetInventoryByProductID(c *gin.Context) {
	id, ok := parsePositiveProductID(c, "product_id")
	if !ok {
		return
	}

	inventory, err := service.GetInventoryByProductID(id)
	if err != nil {
		switch {
		case err == service.ErrInventoryNotFound:
			response.Fail(c, http.StatusNotFound, 2005, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 2002, "查询库存失败")
		}
		return
	}

	response.Success(c, inventory)
}
func ListStockLogs(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 2001, "接口暂未实现")
}

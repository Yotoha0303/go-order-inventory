package handler

import (
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"

	"github.com/gin-gonic/gin"
)

func InitInventory(c *gin.Context) {
	var req request.InitInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlerError(c, err, response.CodeProductParameterError, "请求参数错误")
		return
	}

	if err := service.InitInventory(&req); err != nil {
		handlerError(c, err, response.CodeInitInventoryFailed, "初始化库存错误")
		return
	}

	response.Success(c, nil)
}

func AddInventory(c *gin.Context) {
	var req request.AddInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlerError(c, err, response.CodeParameterError, "请求参数错误")
		return
	}
	if err := service.AddInventory(req); err != nil {
		handlerError(c, err, response.CodeAddInventoryError, "添加库存失败")
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
		handlerError(c, err, response.CodeInventoryNotFound, "查询库存失败")
		return
	}

	response.Success(c, inventory)
}

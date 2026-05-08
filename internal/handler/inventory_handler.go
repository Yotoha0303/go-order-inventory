package handler

import (
	"go-order-inventory/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitInventory(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 2001, "接口暂未实现")
}

func AddInventory(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 2001, "接口暂未实现")
}

func GetInventoryByProductID(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 2001, "接口暂未实现")
}

func ListStockLogs(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 2001, "接口暂未实现")
}

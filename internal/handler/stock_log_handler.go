package handler

import (
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListStockLogs(c *gin.Context) {
	var productID *int64

	productIDStr := c.Query("product_id")
	if productIDStr != "" {
		id, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil || id <= 0 {
			handleError(c, err, http.StatusBadRequest, "无效的产品ID")
			return
		}
		productID = &id
	}

	stockLogs, err := service.ListStockLogsByProductID(productID)
	if err != nil {
		handleError(c, err, response.CodeCreateStockLogFailed, "库存流水日志失败")
		return
	}
	response.Success(c, stockLogs)
}

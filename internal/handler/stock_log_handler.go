package handler

import (
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListStockLogs(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)

	if productID <= 0 && c.Query("product_id") != "" {
		response.Fail(c, http.StatusBadRequest, 2000, "无效的产品ID")
		return
	}

	stockLogs, err := service.ListStockLogsByProductID(&productID)
	if err != nil {
		switch {
		case err == service.ErrStockLogNotFound:
			response.Fail(c, http.StatusNotFound, 2002, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 2001, "查询库存日志失败")
		}
		return
	}
	response.Success(c, stockLogs)
}

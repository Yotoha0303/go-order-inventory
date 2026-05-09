package handler

import (
	"errors"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 1001, "参数错误")
		return
	}

	product, err := service.CreateProduct(req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProductName), errors.Is(err, service.ErrInvalidProductPrice), errors.Is(err, service.ErrInvalidProductDescription):
			response.Fail(c, http.StatusBadRequest, 1002, err.Error())
			return
		default:
			response.Fail(c, 500, 1003, "创建商品失败")
		}
		return
	}

	response.Success(c, product)
}

func ListProducts(c *gin.Context) {

	products, err := service.ListProducts()

	if err != nil {
		response.Fail(c, 500, 1003, err.Error())
		return
	}

	response.Success(c, products)
}

func parsePositiveID(c *gin.Context, paramName string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(paramName), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, 400, 1001, "无效的商品ID")
		return 0, false
	}
	return id, true
}

func GetProductByID(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	product, err := service.GetProductByID(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, http.StatusNotFound, 1002, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, 1003, "查询商品失败")
		}
		return
	}
	response.Success(c, product)
}

func OnSaleProduct(c *gin.Context) {

	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	if err := service.OnSaleProduct(id); err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, 404, 1001, err.Error())
		case errors.Is(err, service.ErrProductOnSaleFailed):
			response.Fail(c, 405, 1002, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 1003, err.Error())
		}
		return
	}
	response.Success(c, nil)
}

func OffSaleProduct(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	if err := service.OffSaleProduct(id); err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, 404, 1001, err.Error())
		case errors.Is(err, service.ErrProductOffSaleFailed):
			response.Fail(c, 405, 1003, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, 1004, err.Error())
		}
		return
	}
	response.Success(c, nil)
}

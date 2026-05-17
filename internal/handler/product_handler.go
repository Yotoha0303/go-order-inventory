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
		response.Fail(c, http.StatusBadRequest, response.CodeParameterError, "参数错误")
		return
	}

	product, err := service.CreateProduct(req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProductName), errors.Is(err, service.ErrInvalidProductPrice), errors.Is(err, service.ErrInvalidProductDescription):
			response.Fail(c, http.StatusBadRequest, response.CodeProductParameterError, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, response.CodeCreateProductFailed, "创建商品失败")
		}
		return
	}

	response.Success(c, product)
}

func ListProducts(c *gin.Context) {

	products, err := service.ListProducts()

	if err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeProductNotFound, err.Error())
		return
	}

	response.Success(c, products)
}

func parsePositiveID(c *gin.Context, paramName string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(paramName), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParameterError, service.ErrInvalidProductID.Error())
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
			response.Fail(c, http.StatusNotFound, response.CodeProductNotFound, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, response.CodeProductNotFound, "查询商品失败")
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
			response.Fail(c, http.StatusNotFound, response.CodeProductNotFound, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, response.CodeProductOnsaleFailed, err.Error())
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
			response.Fail(c, http.StatusNotFound, response.CodeProductNotFound, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, response.CodeProductOffsaleFailed, err.Error())
		}
		return
	}
	response.Success(c, nil)
}

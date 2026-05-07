package handler

import (
	"errors"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"
	"go-order-inventory/internal/service"
	"net/http"

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
			response.Fail(c, http.StatusInternalServerError, 1003, "创建商品失败")
		}
		return
	}

	response.Success(c, product)
}

func ListProducts(c *gin.Context) {

}

func GetProductByID(c *gin.Context) {

}

func OnSaleProduct(c *gin.Context) {

}

func OffSaleProduct(c *gin.Context) {

}

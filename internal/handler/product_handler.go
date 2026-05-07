package handler

import (
	"errors"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/response"

	"github.com/gin-gonic/gin"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

func CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 1001, "参数错误")
		return
	}

	// err := service.CreateProduct(req)
	// if err != nil {
	// 	response.Error(c, 500, 1002, err)
	// 	return
	// }

	response.Success(c, nil)
}

func ListProducts(c *gin.Context) {

}

func GetProductByID(c *gin.Context) {

}

func OnSaleProduct(c *gin.Context) {

}

func OffSaleProduct(c *gin.Context) {

}

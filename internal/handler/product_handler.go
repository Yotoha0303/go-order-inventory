package handler

import (
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
		handleError(c, err, response.CodeCreateProductFailed, "创建商品失败")
		return
	}

	response.Success(c, product)
}

func ListProducts(c *gin.Context) {

	products, err := service.ListProducts()

	if err != nil {
		handleError(c, err, response.CodeQueryProductListFailed, "查询商品列表失败")
		return
	}

	response.Success(c, products)
}

func parsePositiveID(c *gin.Context, paramName string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(paramName), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParameterError, "请求参数错误")
		return 0, false
	}
	return id, true
}

func GetProductByID(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	product, err := service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err, response.CodeQueryProductFailed, "请求商品详情失败")
		return
	}
	response.Success(c, product)
}

func OnSaleProduct(c *gin.Context) {

	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	if err := service.OnSaleProduct(c.Request.Context(), id); err != nil {
		handleError(c, err, response.CodeProductOnSaleFailed, "上架商品失败")
		return
	}
	response.Success(c, nil)
}

func OffSaleProduct(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	if err := service.OffSaleProduct(c.Request.Context(), id); err != nil {
		handleError(c, err, response.CodeProductOffSaleFailed, "下架商品失败")
		return
	}
	response.Success(c, nil)
}

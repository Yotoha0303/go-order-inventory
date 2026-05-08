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
			response.Fail(c, http.StatusInternalServerError, 1003, "创建商品失败")
		}
		return
	}

	response.Success(c, product)
}

func ListProducts(c *gin.Context) {

	products, err := service.ListProducts()
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, 400, 1001, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, 1502, "查询商品列表失败")
		}
		return
	}

	response.Success(c, products)
}

func GetProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 1001, service.ErrInvalidProductID.Error())
		return
	}
	product, err := service.GetProductByID(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, 401, 1002, err.Error())
			return
		default:
			response.Fail(c, 401, 1003, "查询商品失败")
		}
		return
	}
	response.Success(c, product)
}

func OnSaleProduct(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 1001, service.ErrInvalidProductID.Error())
		return
	}
	if err := service.OnSaleProduct(id); err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, 400, 1002, err.Error())
		case errors.Is(err, service.ErrProductOnSaleFailed),
			errors.Is(err, service.ErrProductAlreadyOnSale):
			response.Fail(c, 400, 1001, err.Error())
		default:
			response.Fail(c, 500, 1003, "上架商品失败")
		}
		return
	}
	response.Success(c, nil)
}

func OffSaleProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 1001, service.ErrInvalidProductID.Error())
		return
	}
	if err := service.OffSaleProduct(id); err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			response.Fail(c, 400, 1002, err.Error())
		case errors.Is(err, service.ErrProductOffSaleFailed),
			errors.Is(err, service.ErrProductAlreadyOffSale):
			response.Fail(c, 400, 1001, err.Error())
		default:
			response.Fail(c, 500, 1003, "下架商品失败")
		}
		return
	}
	response.Success(c, nil)
}

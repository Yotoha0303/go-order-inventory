package router

import (
	"go-order-inventory/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()

	registerHealthRouters(r)
	registerAPIRouter(r)
	return r
}

func registerHealthRouters(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})

}

func registerAPIRouter(rg *gin.Engine) {
	apiV1 := rg.Group("/api/v1")

	registerProductAPIRouter(apiV1)
	registerInventoryAPIRouter(apiV1)
}

func registerProductAPIRouter(rg *gin.RouterGroup) {

	rg.POST("/products", handler.CreateProduct)
	rg.GET("/products", handler.ListProducts)
	rg.GET("/products?status=:status", handler.ListProducts)
	rg.GET("/products/:id", handler.GetProductByID)
	rg.PATCH("/products/:id/on-sale", handler.OnSaleProduct)
	rg.PATCH("/products/:id/off-sale", handler.OffSaleProduct)

}

func registerInventoryAPIRouter(rg *gin.RouterGroup) {

	rg.POST("/inventory/init", handler.InitInventory)
	rg.POST("/inventory/add", handler.AddInventory)
	rg.GET("/inventory/products/:product_id", handler.GetInventoryByProductID)
	rg.GET("/stock-logs", handler.ListStockLogs)
}

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
	registerStockLogAPIRouter(apiV1)
	registerOrderAPIRouter(apiV1)
}

func registerProductAPIRouter(rg *gin.RouterGroup) {

	rg.POST("/products", handler.CreateProduct)
	rg.GET("/products", handler.ListProducts)
	rg.GET("/products/:id", handler.GetProductByID)
	rg.PATCH("/products/:id/on-sale", handler.OnSaleProduct)
	rg.PATCH("/products/:id/off-sale", handler.OffSaleProduct)

}

func registerInventoryAPIRouter(rg *gin.RouterGroup) {

	rg.POST("/inventory/init", handler.InitInventory)
	rg.POST("/inventory/add", handler.AddInventory)
	rg.GET("/inventory/products/:product_id", handler.GetInventoryByProductID)
}

func registerStockLogAPIRouter(rg *gin.RouterGroup) {

	rg.GET("/stock-logs", handler.ListStockLogs)

}

func registerOrderAPIRouter(rg *gin.RouterGroup) {

	rg.POST("/orders", handler.CreateOrder)
	rg.GET("/orders/:id", handler.GetOrderByID)
	rg.GET("/orders", handler.ListOrders)

}

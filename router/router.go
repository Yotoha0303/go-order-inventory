package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()

	registerHealthRouters(r)
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
	// apiV1 := rg.Group("/api/v1")
	// registerProductAPIRouter(rg)
	// registerInventoryAPIRouter(rg)
}

func registerProductAPIRouter(rg *gin.Engine) {
	// POST /api/v1/products
	// GET /api/v1/products
	// GET /api/v1/products/:id
	// PATCH /api/v1/products/:id/on-sale
	// PATCH /api/v1/products/:id/off-sale

}

func registerInventoryAPIRouter(rg *gin.Engine) {
	// POST /api/v1/inventory/init
	// POST /api/v1/inventory/add
	// GET /api/v1/inventory/products/:product_id
	// GET /api/v1/stock-logs
}

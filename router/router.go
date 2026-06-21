package router

import (
	"go-order-inventory/internal/handler"
	"go-order-inventory/internal/middleware"
	"go-order-inventory/internal/service"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouters(db *gorm.DB, logger *slog.Logger, timeout time.Duration) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.RequestID(),
		middleware.AccessLog(logger),
		middleware.TimeoutMiddleware(timeout),
		middleware.Recovery(logger),
	)

	productService := service.NewProductService(db)
	inventoryService := service.NewInventoryService(db)
	stockLogService := service.NewStockLogService(db)
	orderService := service.NewOrderService(db)

	productHandler := handler.NewProductHandler(productService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	stockLogHandler := handler.NewStockLogHandler(stockLogService)
	orderHandler := handler.NewOrderHandler(orderService)

	registerHealthRouters(r)
	registerAPIRouter(r, productHandler, inventoryHandler, stockLogHandler, orderHandler)
	return r
}

func registerHealthRouters(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"message": "success",
			},
		})
	})

}

func registerAPIRouter(
	rg *gin.Engine,
	productHandler *handler.ProductHandler,
	inventoryHandler *handler.InventoryHandler,
	stockLogHandler *handler.StockLogHandler,
	orderHandler *handler.OrderHandler,
) {
	apiV1 := rg.Group("/api/v1")

	registerProductAPIRouter(apiV1, productHandler)
	registerInventoryAPIRouter(apiV1, inventoryHandler)
	registerStockLogAPIRouter(apiV1, stockLogHandler)
	registerOrderAPIRouter(apiV1, orderHandler)
}

func registerProductAPIRouter(rg *gin.RouterGroup, productHandler *handler.ProductHandler) {

	rg.POST("/products", productHandler.CreateProduct)
	rg.GET("/products", productHandler.ListProducts)
	rg.GET("/products/:id", productHandler.GetProductByID)
	rg.PATCH("/products/:id/on-sale", productHandler.OnSaleProduct)
	rg.PATCH("/products/:id/off-sale", productHandler.OffSaleProduct)

}

func registerInventoryAPIRouter(rg *gin.RouterGroup, inventoryHandler *handler.InventoryHandler) {

	rg.POST("/inventory/init", inventoryHandler.InitInventory)
	rg.POST("/inventory/add", inventoryHandler.AddInventory)
	rg.GET("/inventory/products/:product_id", inventoryHandler.GetInventoryByProductID)
}

func registerStockLogAPIRouter(rg *gin.RouterGroup, stockLogHandler *handler.StockLogHandler) {

	rg.GET("/stock-logs", stockLogHandler.ListStockLogs)

}

func registerOrderAPIRouter(rg *gin.RouterGroup, orderHandler *handler.OrderHandler) {

	rg.POST("/orders", orderHandler.CreateOrder)
	rg.GET("/orders/:id", orderHandler.GetOrderByID)
	rg.GET("/orders", orderHandler.ListOrders)
	rg.PATCH("/orders/:id/cancel", orderHandler.CancelOrders)
	rg.PATCH("/orders/:id/pay", orderHandler.PayOrder)
	rg.PATCH("/orders/:id/finish", orderHandler.FinishOrder)

}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/product-management-server/internal/handler"
)

// Register mounts all application routes onto the given Gin engine.
func Register(router *gin.Engine, productHandler *handler.ProductHandler) {
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products", productHandler.GetProducts)
	router.GET("/products/:id", productHandler.GetProductByID)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	router.DELETE("/products/:id", productHandler.DeleteProduct)
}

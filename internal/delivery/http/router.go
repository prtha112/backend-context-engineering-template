package http

import (
	"backend-context-engineering-template/internal/delivery/http/handlers"
	"backend-context-engineering-template/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRouter(productHandler *handlers.ProductHandler, logger *logrus.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.ErrorHandler(logger))

	api := r.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("", productHandler.GetProducts)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

	return r
}

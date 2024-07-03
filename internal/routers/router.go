package router

import (
	"basic-trade-app/internal/controllers"
	"basic-trade-app/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {

	router := gin.Default()

	// router grouping
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", controllers.AdminRegister)
		authRouter.POST("/login", controllers.AdminLogin)
	}

	productRouter := router.Group("/products")
	{
		// public API
		productRouter.GET("/", controllers.GetProducts)
		productRouter.GET("/:productUUID", controllers.GetProduct)

		// set layer authentication jika perlu authentication di request (bearer token)
		productRouter.Use(middlewares.Authentication()) 
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.DELETE("/:productUUID", middlewares.ProductAuthorization(), controllers.DeleteProduct)
		productRouter.PUT("/:productUUID", middlewares.ProductAuthorization(), controllers.UpdateProduct)
		
		// variant routes
		productRouter.POST("/variants", middlewares.ProductAuthorization(), controllers.CreateVariant)
		productRouter.GET("/variants", controllers.GetVariants)
		productRouter.GET("/variants/:variantUUID", controllers.GetVariant)
	}
	return router
}
package router

import (
	"basic-trade-app/internal/controllers"

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

	return router

}
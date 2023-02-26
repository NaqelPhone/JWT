package routes

import (
	"system/controllers"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/login", controllers.GetLogin())
	incomingRoutes.POST("/login", controllers.Login())
}

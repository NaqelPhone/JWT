package routes

import (
	"net/http"
	"system/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/signup", controllers.GetSignUp())
	incomingRoutes.POST("/signup", controllers.SignUp())
	incomingRoutes.GET("/admin", func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusOK, gin.H{"Logged in": "Admin"})
	})
}

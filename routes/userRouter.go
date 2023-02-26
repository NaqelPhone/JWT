package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusOK, gin.H{"Logged in": "User"})
	})
}

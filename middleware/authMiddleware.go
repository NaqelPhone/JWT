package middleware

import (
	"net/http"
	"system/helpers"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		clientToken, _ := ginCtx.Cookie("token")

		if clientToken == "" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Please login"})
			ginCtx.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			ginCtx.Abort()
			return
		}

		ginCtx.Set("username", claims.UserName)
		ginCtx.Set("firstname", claims.FirstName)
		ginCtx.Set("lastname", claims.LastName)
		ginCtx.Set("userId", claims.UID)
		ginCtx.Set("role", claims.Role)
		ginCtx.Set("hospital", claims.Hospital)

		ginCtx.Next()
	}

}

func Admin() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		clientToken, _ := ginCtx.Cookie("token")

		if clientToken == "" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Please login"})
			ginCtx.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			ginCtx.Abort()
			return
		}

		if claims.Role != "admin" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "not admin"})
			ginCtx.Abort()
			return
		}

		ginCtx.Set("username", claims.UserName)
		ginCtx.Set("firstname", claims.FirstName)
		ginCtx.Set("lastname", claims.LastName)
		ginCtx.Set("userId", claims.UID)
		ginCtx.Set("role", claims.Role)
		ginCtx.Set("hospital", claims.Hospital)

		ginCtx.Next()
	}
}

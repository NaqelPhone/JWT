package controllers

import (
	"context"
	"net/http"
	helper "system/helpers"
	"system/models"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
)

func Login() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var foundUser models.User
		UserName := ginCtx.Request.FormValue("username")
		Password := ginCtx.Request.FormValue("password")
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"username": UserName}).Decode(&foundUser)

		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		passwordIsValid, message := VerifyPassword(Password, foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(foundUser.UserName, foundUser.FirstName, foundUser.LastName, foundUser.UserID, foundUser.Role, foundUser.Hospital)
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		ginCtx.Writer.Header().Set("Token", token)
		ginCtx.Writer.Header().Set("Authorization", "Bearer "+token)
		ginCtx.SetCookie("token", token, 36000, "/", "localhost", false, true)
		ginCtx.JSON(http.StatusOK, gin.H{"message": "Logged in"})

	}
}

func GetLogin() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.HTML(http.StatusOK, "login.html", nil)
	}
}

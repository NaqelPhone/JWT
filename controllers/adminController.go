package controllers

import (
	"context"
	"log"
	"net/http"
	helper "system/helpers"
	"system/models"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		user.FirstName = ginCtx.Request.FormValue("firstname")
		user.LastName = ginCtx.Request.FormValue("lastname")
		user.Hospital = ginCtx.Request.FormValue("hospital")
		user.Email = ginCtx.Request.FormValue("email")
		user.Phone = ginCtx.Request.FormValue("phone")
		user.UserName = ginCtx.Request.FormValue("username")
		user.Password = ginCtx.Request.FormValue("password")
		user.Role = ginCtx.Request.FormValue("role")

		defer cancel()

		validationError := validate.Struct(user)
		if validationError != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.UserName})

		if err != nil {
			log.Panic(err)
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "counting document error"})
			return
		}
		if count > 0 {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "username or phone is already in use"})
			return
		}

		password := HashPassword(user.Password)
		user.Password = password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "counting document error"})
			return
		}

		if count > 0 {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "username or phone is already in use"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(user.UserName, user.FirstName, user.LastName, user.UserID, user.Role, user.Hospital)
		user.Token = token
		user.RefreshToken = refreshToken
		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			message := insertErr.Error()
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": message})
			return
		}
		defer cancel()
		ginCtx.JSON(http.StatusOK, "User inserted successfully")
	}
}

func GetSignUp() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.HTML(http.StatusOK, "signup.html", nil)
	}
}

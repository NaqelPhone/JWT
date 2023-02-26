package controllers

import (
	"log"
	"system/database"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, inputPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(inputPassword), []byte(userPassword))
	check := true
	var message = "Logging in... "
	if err != nil {
		message = err.Error()
		check = false
	}
	return check, message
}

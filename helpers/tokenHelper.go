package helpers

import (
	"context"
	"fmt"
	"log"
	"system/database"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	UserName  string
	FirstName string
	LastName  string
	UID       string
	Role      string
	Hospital  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string = "VeRySeCrEtKeY"

func GenerateAllTokens(username string, firstname string, lastname string, uid string, role string, hospital string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		UserName:  username,
		FirstName: firstname,
		LastName:  lastname,
		UID:       uid,
		Role:      role,
		Hospital:  hospital,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, message string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		message = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		message = fmt.Sprintf("invalid token")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		message = fmt.Sprintf("expired token")
		return
	}

	return claims, message
}

func UpdateAllTokens(signedToken, signedRefreshToken, userID string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObject primitive.D

	updateObject = append(updateObject, bson.E{Key: "token", Value: signedToken})
	updateObject = append(updateObject, bson.E{Key: "refreshToken", Value: signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updatedAt", Value: UpdatedAt})

	upsert := true
	filter := bson.M{"userId": userID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObject},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}

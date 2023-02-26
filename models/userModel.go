package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"firstname" validate:"required"`
	LastName     string             `json:"lastname" validate:"required"`
	UserName     string             `json:"username" validate:"required"`
	Password     string             `json:"password" validate:"required"`
	UserID       string             `json:"userId"`
	Role         string             `json:"role" validate:"required"`
	Hospital     string             `json:"hospital" validate:"required"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone" validate:"required"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}

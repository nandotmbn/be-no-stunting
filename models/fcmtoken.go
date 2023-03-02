package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FCMToken struct {
	FCMToken string             `json:"fcmtoken" validate:"required"`
	UserId   primitive.ObjectID `json:"userid,omitempty"`
}

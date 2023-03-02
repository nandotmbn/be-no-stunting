package views

import "go.mongodb.org/mongo-driver/bson/primitive"

type FCMTokenView struct {
	Id       primitive.ObjectID `json:"_id,omitempty" validate:"required"`
	FCMToken string             `json:"fcmtoken" validate:"required"`
	UserId   primitive.ObjectID `json:"userid,omitempty"`
}

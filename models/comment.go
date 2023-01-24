package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	PostId    primitive.ObjectID `json:"postId,omitempty" bson:"postid,omitempty"`
	UserId    primitive.ObjectID `json:"userId,omitempty" bson:"userid,omitempty"`
	Content   string             `json:"content,omitempty" validate:"required,min=0"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

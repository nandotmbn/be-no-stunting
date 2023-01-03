package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	PostId    primitive.ObjectID `json:"postid,omitempty"`
	Content   string             `json:"content,omitempty" validate:"required,min=0"`
	CreatedAt time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

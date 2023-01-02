package views

import "go.mongodb.org/mongo-driver/bson/primitive"

type RolesWithId struct {
	Name        string             `json:"name" validate:"required,gte=1,lte=255"`
	DisplayName string             `json:"displayName" validate:"required,gte=1,lte=255"`
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

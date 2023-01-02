package views

import "go.mongodb.org/mongo-driver/bson/primitive"

type FacilityMeasureFindGet struct {
	FirstName  string             `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	LastName   string             `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

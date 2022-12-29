package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	Height     int                `json:"height,omitempty" validate:"required,min=0"`
	Weight     int                `json:"weight,omitempty" validate:"required,min=0"`
	HeartRate  int                `json:"heartrate,omitempty" validate:"required,min=0"`
	Temp       int                `json:"temp,omitempty" validate:"required,min=0"`
	ChildId    primitive.ObjectID `json:"childid,omitempty" validate:"required"`
	FacilityId primitive.ObjectID `json:"facilityid,omitempty"`
	CreatedAt  time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

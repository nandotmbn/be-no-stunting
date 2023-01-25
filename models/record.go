package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	IsChecked  bool               `json:"ischecked"`
	Height     int                `json:"height,omitempty" validate:"required,min=0"`
	Weight     int                `json:"weight,omitempty" validate:"required,min=0"`
	PatientId  primitive.ObjectID `json:"patientid,omitempty" validate:"required"`
	FacilityId primitive.ObjectID `json:"facilityid,omitempty"`
	CreatedAt  time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

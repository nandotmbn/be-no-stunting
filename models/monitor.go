package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Monitor struct {
	IsChecked     bool               `json:"ischecked"`
	PatientTypeId primitive.ObjectID `json:"patienttypeid,omitempty"`
	PatientId     primitive.ObjectID `json:"patientid,omitempty"`
	FacilityId    primitive.ObjectID `json:"facilityid,omitempty"`
	Content       string             `json:"content,omitempty" validate:"required,min=0"`
	CreatedAt     time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

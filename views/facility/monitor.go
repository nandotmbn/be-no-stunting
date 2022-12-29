package views

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	FirstName string    `json:"firstname,omitempty" validate:"required,min=3,max=255"`
	LastName  string    `json:"lastname,omitempty" validate:"required,min=3,max=255"`
	Address   string    `json:"address,omitempty" validate:"required,min=3,max=255"`
	Id        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UpdatedAt time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
}

type FacilityMonitorGet struct {
	Id            string             `json:"_id,omitempty" bson:"_id,omitempty"`
	IsChecked     bool               `json:"ischecked"`
	PatientTypeId primitive.ObjectID `json:"patienttypeid,omitempty"`
	PatientId     primitive.ObjectID `json:"patientid,omitempty"`
	Patient       []User             `json:"patient,omitempty"`
	FacilityId    primitive.ObjectID `json:"facilityid,omitempty"`
	Content       string             `json:"content,omitempty" validate:"required,min=0"`
	CreatedAt     time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

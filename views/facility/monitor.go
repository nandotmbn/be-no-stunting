package views

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	FirstName string    `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	LastName  string    `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	Address   string    `json:"address,omitempty" validate:"required,min=3,max=255"`
	Id        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedat,omitempty"`
}

type FacilityMonitorGet struct {
	Id            string             `json:"_id,omitempty" bson:"_id,omitempty"`
	IsChecked     bool               `json:"isChecked"`
	PatientTypeId primitive.ObjectID `json:"patientTypeId,omitempty"`
	PatientId     primitive.ObjectID `json:"patientId,omitempty"`
	Patient       []User             `json:"patient,omitempty"`
	FacilityId    primitive.ObjectID `json:"facilityId,omitempty"`
	Content       string             `json:"content,omitempty" validate:"required,min=0"`
	CreatedAt     time.Time          `json:"createdAt,omitempty" bson:"createdat,omitempty"`
}

type FacilityMonitorFindByIdFinal struct {
	Patient    []User_            `json:"patient,omitempty"`
	PatientId  primitive.ObjectID `json:"patientId,omitempty" bson:"patientid,omitempty"`
	Height     int                `json:"height,omitempty" validate:"required,min=0"`
	Weight     int                `json:"weight,omitempty" validate:"required,min=0"`
	HeartRate  int                `json:"heartRate,omitempty" validate:"required,min=0"`
	Temp       int                `json:"temp,omitempty" validate:"required,min=0"`
	FacilityId primitive.ObjectID `json:"facilityId,omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"createdat,omitempty"`
	IsChecked  bool               `json:"isChecked"`
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content    string             `json:"content,omitempty" validate:"required,min=0"`
}

package views

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User_ struct {
	FirstName string `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	LastName  string `json:"lastName,omitempty" validate:"required,min=3,max=255"`
}

type FacilityMeasureFindGet struct {
	FirstName  string             `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	LastName   string             `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" validate:"required,min=3,max=255"`
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type FacilityMeasureGet struct {
	Id         string             `json:"_id,omitempty" bson:"_id,omitempty"`
	IsChecked  bool               `json:"isChecked"`
	Height     int                `json:"height,omitempty" validate:"required,min=0"`
	Weight     int                `json:"weight,omitempty" validate:"required,min=0"`
	HeartRate  int                `json:"heartRate,omitempty" validate:"required,min=0"`
	Temp       int                `json:"temp,omitempty" validate:"required,min=0"`
	PatientId  primitive.ObjectID `json:"patientId,omitempty"`
	Patient    []User             `json:"patient,omitempty"`
	FacilityId primitive.ObjectID `json:"facilityId,omitempty"`
	Content    string             `json:"content,omitempty" validate:"required,min=0"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"createdat,omitempty"`
}

type FacilityMonitorFindById struct {
	IsChecked bool               `json:"isChecked"`
	CreatedAt time.Time          `json:"createdAt,omitempty" validate:"required,min=3,max=255"`
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PatientId primitive.ObjectID `json:"patientId,omitempty" bson:"patientid,omitempty"`
	Content   string             `json:"content,omitempty" validate:"required,min=0"`
}

type FacilityMonitorFindByIdView struct {
	IsChecked bool               `json:"isChecked"`
	CreatedAt time.Time          `json:"createdAt,omitempty" validate:"required,min=3,max=255"`
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PatientId primitive.ObjectID `json:"patientId,omitempty" bson:"patientid,omitempty"`
	Patient   []User_            `json:"patient,omitempty"`
	Content   string             `json:"content,omitempty" validate:"required,min=0"`
}

type FacilityRecordFindById struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IsChecked  bool               `json:"isChecked"`
	Patient    []User_            `json:"patient,omitempty"`
	PatientId  primitive.ObjectID `json:"patientId,omitempty" bson:"patientid,omitempty"`
	Height     int                `json:"height,omitempty" validate:"required,min=0"`
	Weight     int                `json:"weight,omitempty" validate:"required,min=0"`
	HeartRate  int                `json:"heartRate,omitempty" validate:"required,min=0"`
	Temp       int                `json:"temp,omitempty" validate:"required,min=0"`
	FacilityId primitive.ObjectID `json:"facilityId,omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"createdat,omitempty"`
}

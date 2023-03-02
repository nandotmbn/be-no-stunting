package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	FirstName  string             `json:"firstname,omitempty" validate:"required,min=3,max=255"`
	LastName   string             `json:"lastname,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	Password   string             `json:"password,omitempty" validate:"required,min=3,max=255"`
	Address    string             `json:"address,omitempty" validate:"required,min=3,max=255"`
	RolesId    primitive.ObjectID `json:"rolesid,omitempty" validate:"required"`
	ParentId   primitive.ObjectID `json:"parentid,omitempty"`
	Id         string             `json:"id,omitempty" bson:"_id,omitempty"`
	UpdatedAt  time.Time          `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	CreatedAt  time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

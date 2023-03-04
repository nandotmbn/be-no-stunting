package views

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserViewResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type InputLogin struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	FCMToken   string `json:"fcmtoken" validate:"required"`
}

type InputLogout struct {
	FCMToken string `json:"fcmtoken" validate:"required"`
	UserId   string `json:"userid" validate:"required"`
}

type UserNoPassword struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName  string             `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	LastName   string             `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	ParentId   primitive.ObjectID `json:"parentId,omitempty"`
	IsMale     bool               `json:"isMale,omitempty"`
	Address    string             `json:"address,omitempty" validate:"required,min=3,max=255"`
	RolesId    primitive.ObjectID `json:"rolesId,omitempty" validate:"required"`
	BornAt     time.Time          `json:"bornAt,omitempty" bson:"bornat,omitempty"`
	UpdatedAt  time.Time          `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
}

type UserUpdate struct {
	FirstName  string             `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	LastName   string             `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	ParentId   primitive.ObjectID `json:"parentId,omitempty"`
	Address    string             `json:"address,omitempty" validate:"required,min=3,max=255"`
	RolesId    primitive.ObjectID `json:"rolesId,omitempty" validate:"required"`
}

type UserOnlyId struct {
	Id string `json:"_id,omitempty" bson:"_id,omitempty"`
}

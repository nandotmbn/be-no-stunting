package views

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateUserViewResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type InputLogin struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type UserNoPassword struct {
	FirstName  string             `json:"firstname,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	LastName   string             `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	ParentId   primitive.ObjectID `json:"parentid,omitempty"`
	RolesId    primitive.ObjectID `json:"rolesid,omitempty" validate:"required"`
	Id         string             `json:"_id,omitempty" bson:"_id,omitempty"`
}

type UserOnlyId struct {
	Id string `json:"_id,omitempty" bson:"_id,omitempty"`
}

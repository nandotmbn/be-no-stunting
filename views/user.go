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
	FirstName  string             `json:"firstName,omitempty" validate:"required,min=3,max=255"`
	Identifier string             `json:"identifier,omitempty" validate:"required,min=3,max=55"`
	LastName   string             `json:"lastName,omitempty" validate:"required,min=3,max=255"`
	ParentId   primitive.ObjectID `json:"parentId,omitempty"`
	RolesId    primitive.ObjectID `json:"rolesId,omitempty" validate:"required"`
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type UserOnlyId struct {
	Id string `json:"_id,omitempty" bson:"_id,omitempty"`
}

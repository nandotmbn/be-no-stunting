package models

type Roles struct {
	Name        string `json:"name" validate:"required,gte=1,lte=255"`
	DisplayName string `json:"displayname" validate:"required,gte=1,lte=255"`
}

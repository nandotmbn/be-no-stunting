package views

type Master struct {
	Roles   []string `json:"roles" validate:"required,gte=1,lte=1,dive"`
	MomType []string `json:"momType" validate:"required,gte=1,lte=1,dive"`
}

type MasterResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

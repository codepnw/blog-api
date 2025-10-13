package categoryhandler

type CategoryCreateReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type CategoryUpdateReq struct {
	Name        *string `json:"name,omitempty" validate:"omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
}

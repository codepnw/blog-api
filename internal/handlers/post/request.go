package posthandler

type PostCreate struct {
	Title      string  `json:"title" validate:"required"`
	Content    *string `json:"content,omitempty" validate:"omitempty"`
	CategoryID *string `json:"category_id,omitempty" validate:"omitempty"`
}

type PostUpdate struct {
	Title      *string `json:"title,omitempty" validate:"omitempty"`
	Content    *string `json:"content,omitempty" validate:"omitempty"`
	CategoryID *string `json:"category_id,omitempty" validate:"omitempty"`
}

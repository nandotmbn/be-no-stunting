package views

type CommentInput struct {
	Content string `json:"content" validate:"required,min=1"`
}

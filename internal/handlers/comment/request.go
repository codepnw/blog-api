package commenthandler

type CommentReq struct {
	Content string `json:"content" validate:"required"`
}

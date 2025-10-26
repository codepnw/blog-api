package commenthandler

import (
	"strconv"

	commentdomain "github.com/codepnw/blog-api/internal/domains/comment"
	"github.com/codepnw/blog-api/internal/handlers"
	"github.com/codepnw/blog-api/internal/middleware"
	commentusecase "github.com/codepnw/blog-api/internal/usecases/comment"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc commentusecase.Usecase
}

func NewCommentHandler(uc commentusecase.Usecase) *handler {
	return &handler{uc: uc}
}

// Create Comment
// @Summary Create Comment
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID"
// @Param data body commenthandler.CommentReq true "New comment"
// @Success 201 {object} commentdomain.Comment
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts{post_id}/comments [post]
func (h *handler) CreateComment(ctx *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		return handlers.Unauthorized(ctx, err.Error())
	}

	postID := ctx.Params(handlers.ParamKeyPostID)

	req := new(CommentReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &commentdomain.Comment{
		UserID:  user.UserID,
		PostID:  postID,
		Content: req.Content,
	}
	result, err := h.uc.CreateComment(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.Created(ctx, result)
}

// Get Comment By Post
// @Summary Get Comment By Post
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {object} []commentdomain.Comment
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts{post_id}/comments [get]
func (h *handler) GetCommentByPost(ctx *fiber.Ctx) error {
	postID := ctx.Params(handlers.ParamKeyPostID)

	result, err := h.uc.GetCommentByPost(ctx.Context(), postID)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.Success(ctx, result)
}

// Edit Comment
// @Summary Edit Comment
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID"
// @Param comment_id path string true "Comment ID"
// @Param data body commenthandler.CommentReq true "New comment"
// @Success 200 {object} commentdomain.Comment
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts{post_id}/comments/{comment_id} [patch]
func (h *handler) EditComment(ctx *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		return handlers.Unauthorized(ctx, err.Error())
	}

	postID := ctx.Params(handlers.ParamKeyPostID)
	commentID, _ := ctx.ParamsInt(handlers.ParamKeyCommentID)

	req := new(CommentReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &commentdomain.Comment{
		ID:      int64(commentID),
		UserID:  user.UserID,
		PostID:  postID,
		Content: req.Content,
	}
	result, err := h.uc.EditComment(ctx.Context(), input, user.Role)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.Success(ctx, result)
}

// Delete Comment
// @Summary Delete Comment
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID"
// @Param comment_id path string true "Comment ID"
// @Success 204 {object} handlers.EmptyRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts{post_id}/comments/{comment_id} [delete]
func (h *handler) DeleteComment(ctx *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		return handlers.Unauthorized(ctx, err.Error())
	}

	commentID := ctx.Params(handlers.ParamKeyCommentID)
	id, _ := strconv.ParseInt(commentID, 10, 64)

	if err = h.uc.DeleteComment(ctx.Context(), id, user.UserID, user.Role); err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.NoContent(ctx)
}

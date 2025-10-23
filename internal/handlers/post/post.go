package posthandler

import (
	"errors"

	postdomain "github.com/codepnw/blog-api/internal/domains/post"
	"github.com/codepnw/blog-api/internal/handlers"
	"github.com/codepnw/blog-api/internal/middleware"
	postusecase "github.com/codepnw/blog-api/internal/usecases/post"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
	"github.com/codepnw/blog-api/internal/utils/errs"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc postusecase.Usecase
}

func NewPostHandler(uc postusecase.Usecase) *handler {
	return &handler{uc: uc}
}

// Create Post
// @Summary Create Post
// @Tags posts
// @Accept json
// @Produce json
// @Param data body posthandler.PostCreateReq true "Post data"
// @Success 201 {object} postdomain.Post
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts [post]
func (h *handler) Create(ctx *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		return handlers.Unauthorized(ctx, err.Error())
	}

	req := new(PostCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &postdomain.Post{
		AuthorID:   user.UserID,
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: &req.CategoryID,
	}

	result, err := h.uc.Create(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Created(ctx, result)
}

// Get Post By ID
// @Summary Get Post By ID
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} postdomain.Post
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts/{post_id} [get]
func (h *handler) GetByID(ctx *fiber.Ctx) error {
	postID := ctx.Params(handlers.ParamKeyPostID)

	result, err := h.uc.GetByID(ctx.Context(), postID)
	if err != nil {
		if errors.Is(err, errs.ErrPostNotFound) {
			return handlers.NotFound(ctx, err.Error())
		}
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Get Post By User
// @Summary Get Post By User
// @Tags posts
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} []postdomain.Post
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users/{user_id}/posts [get]
func (h *handler) GetByUserID(ctx *fiber.Ctx) error {
	authorID := ctx.Params(handlers.ParamKeyAuthorID)

	result, err := h.uc.GetByAuthorID(ctx.Context(), authorID)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Get Posts
// @Summary Get Posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} []postdomain.Post
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts [get]
func (h *handler) GetAll(ctx *fiber.Ctx) error {
	result, err := h.uc.GetAll(ctx.Context())
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Edit Post
// @Summary Edit Post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param data body posthandler.PostUpdateReq true "Post data"
// @Success 200 {object} postdomain.Post
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts/{post_id} [patch]
func (h *handler) Update(ctx *fiber.Ctx) error {
	postID := ctx.Params(handlers.ParamKeyPostID)

	req := new(PostUpdateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	ok, err := h.checkPermissions(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserUnauthorized):
			return handlers.Unauthorized(ctx, err.Error())
		case errors.Is(err, errs.ErrPostNotFound):
			return handlers.NotFound(ctx, err.Error())
		default:
			return handlers.InternalServerError(ctx, err)
		}
	}
	if !ok {
		return handlers.Forbidden(ctx, "no permissions")
	}

	input := h.validateUpdate(postID, req)
	result, err := h.uc.Update(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Delete Post
// @Summary Delete Post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 204 {object} handlers.EmptyRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 403 {object} handlers.ForbiddenRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /posts/{post_id} [delete]
func (h *handler) Delete(ctx *fiber.Ctx) error {
	postID := ctx.Params(handlers.ParamKeyPostID)

	ok, err := h.checkPermissions(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserUnauthorized):
			return handlers.Unauthorized(ctx, err.Error())
		case errors.Is(err, errs.ErrPostNotFound):
			return handlers.NotFound(ctx, err.Error())
		default:
			return handlers.InternalServerError(ctx, err)
		}
	}
	if !ok {
		return handlers.Forbidden(ctx, "no permissions")
	}

	if err := h.uc.Delete(ctx.Context(), postID); err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.NoContent(ctx)
}

func (h *handler) validateUpdate(postID string, req *PostUpdateReq) *postdomain.Post {
	newPost := new(postdomain.Post)
	if req.Title != nil {
		newPost.Title = *req.Title
	}
	if req.Content != nil {
		newPost.Content = *req.Content
	}
	if req.CategoryID != nil {
		newPost.CategoryID = req.CategoryID
	}
	newPost.ID = postID

	return newPost
}

func (h *handler) checkPermissions(ctx *fiber.Ctx, postID string) (bool, error) {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		return false, errs.ErrUserUnauthorized
	}

	post, err := h.uc.GetByID(ctx.Context(), postID)
	if err != nil {
		return false, errs.ErrPostNotFound
	}

	if post.AuthorID == user.UserID || user.Role == string(userusecase.RoleAdmin) {
		return true, nil
	}
	return false, nil
}

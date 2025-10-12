package posthandler

import (
	postdomain "github.com/codepnw/blog-api/internal/domains/post"
	"github.com/codepnw/blog-api/internal/handlers"
	postusecase "github.com/codepnw/blog-api/internal/usecases/post"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc postusecase.Usecase
}

func NewPostHandler(uc postusecase.Usecase) *handler {
	return &handler{uc: uc}
}

func (h *handler) Create(ctx *fiber.Ctx) error {
	req := new(PostCreate)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &postdomain.Post{
		// TODO: Get AuthorID From Context Later
		AuthorID:   "",
		Title:      req.Title,
		Content:    *req.Content,
		CategoryID: *req.CategoryID,
	}

	result, err := h.uc.Create(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Created(ctx, result)
}

func (h *handler) GetByID(ctx *fiber.Ctx) error {
	postID := ctx.Params("post_id")

	result, err := h.uc.GetByID(ctx.Context(), postID)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) GetByAuthorID(ctx *fiber.Ctx) error {
	authorID := ctx.Query("author_id")

	result, err := h.uc.GetByAuthorID(ctx.Context(), authorID)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) GetAll(ctx *fiber.Ctx) error {
	result, err := h.uc.GetAll(ctx.Context())
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) Update(ctx *fiber.Ctx) error {
	postID := ctx.Params("post_id")

	req := new(PostUpdate)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &postdomain.Post{
		ID:         postID,
		Title:      *req.Title,
		Content:    *req.Content,
		CategoryID: *req.CategoryID,
	}
	result, err := h.uc.Update(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) Delete(ctx *fiber.Ctx) error {
	postID := ctx.Params("post_id")

	if err := h.uc.Delete(ctx.Context(), postID); err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.NoContent(ctx)
}

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
	req := new(PostCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &postdomain.Post{
		// TODO: Get AuthorID From Context Later
		AuthorID:   "a9bbad43-51e5-424c-8ec5-5f40cd89b701",
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

func (h *handler) GetByID(ctx *fiber.Ctx) error {
	postID := ctx.Params(handlers.ParamKeyPostID)

	result, err := h.uc.GetByID(ctx.Context(), postID)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) GetByAuthorID(ctx *fiber.Ctx) error {
	authorID := ctx.Params(handlers.ParamKeyAuthorID)

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
	postID := ctx.Params(handlers.ParamKeyPostID)

	req := new(PostUpdateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := new(postdomain.Post)
	if req.Title != nil {
		input.Title = *req.Title
	}
	if req.Content != nil {
		input.Content = *req.Content
	}
	if req.CategoryID != nil {
		input.CategoryID = req.CategoryID
	}
	input.ID = postID

	result, err := h.uc.Update(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) Delete(ctx *fiber.Ctx) error {
	postID := ctx.Params(handlers.ParamKeyPostID)

	if err := h.uc.Delete(ctx.Context(), postID); err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.NoContent(ctx)
}

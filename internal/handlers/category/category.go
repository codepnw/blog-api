package categoryhandler

import (
	categorydomain "github.com/codepnw/blog-api/internal/domains/category"
	"github.com/codepnw/blog-api/internal/handlers"
	categoryusecase "github.com/codepnw/blog-api/internal/usecases/category"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc categoryusecase.Usecase
}

func NewCategoryHandler(uc categoryusecase.Usecase) *handler {
	return &handler{uc: uc}
}

func (h *handler) Create(ctx *fiber.Ctx) error {
	req := new(CategoryCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &categorydomain.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.uc.Create(ctx.Context(), input); err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Created(ctx, "added new category")
}

func (h *handler) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyCategoryID)

	result, err := h.uc.GetByID(ctx.Context(), id)
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
	id := ctx.Params(handlers.ParamKeyCategoryID)

	req := new(CategoryUpdateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := new(categorydomain.Category)
	if req.Name != nil {
		input.Name = *req.Name
	}
	if req.Description != nil {
		input.Description = *req.Description
	}
	input.ID = id

	if err := h.uc.Update(ctx.Context(), input); err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, "category updated")
}

func (h *handler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyCategoryID)

	if err := h.uc.Delete(ctx.Context(), id); err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.NoContent(ctx)
}

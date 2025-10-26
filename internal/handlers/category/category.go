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

// Create Category
// @Summary Create Category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body categoryhandler.CategoryCreateReq true "New category"
// @Success 201 {object} categorydomain.Category
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /categories [post]
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

// Get Category
// @Summary Get Category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Success 200 {object} categorydomain.Category
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /categories/category_id [get]
func (h *handler) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyCategoryID)

	result, err := h.uc.GetByID(ctx.Context(), id)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.Success(ctx, result)
}

// Create Categories
// @Summary Create Categories
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []categorydomain.Category
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /categories [get]
func (h *handler) GetAll(ctx *fiber.Ctx) error {
	result, err := h.uc.GetAll(ctx.Context())
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.Success(ctx, result)
}

// Update Category
// @Summary Update Category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Param data body categoryhandler.CategoryUpdateReq true "New category"
// @Success 200 {object} categorydomain.Category
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /categories/category_id [patch]
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

// Delete Category
// @Summary Delete Category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Success 204 {object} handlers.EmptyRes
// @Failure 401 {object} handlers.UnauthorizedRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /categories/category_id [delete]
func (h *handler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyCategoryID)

	if err := h.uc.Delete(ctx.Context(), id); err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.NoContent(ctx)
}

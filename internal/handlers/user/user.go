package userhandler

import (
	userdomain "github.com/codepnw/blog-api/internal/domains/user"
	"github.com/codepnw/blog-api/internal/handlers"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc userusecase.Usecase
}

func NewUserHandler(uc userusecase.Usecase) *handler {
	return &handler{uc: uc}
}

func (h *handler) CreateUser(ctx *fiber.Ctx) error {
	req := new(UserCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &userdomain.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password,
	}
	result, err := h.uc.CreateUser(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Created(ctx, result)
}

func (h *handler) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyUserID)

	result, err := h.uc.GetUser(ctx.Context(), id)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) GetAllUsers(ctx *fiber.Ctx) error {
	result, err := h.uc.GetAllUsers(ctx.Context())
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyUserID)

	req := new(UserUpdateReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := new(userdomain.User)
	if req.FirstName != nil {
		input.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		input.LastName = *req.LastName
	}
	input.ID = id

	result, err := h.uc.UpdateUser(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

func (h *handler) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyUserID)

	if err := h.uc.DeleteUser(ctx.Context(), id); err != nil {
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.NoContent(ctx)
}

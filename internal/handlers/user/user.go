package userhandler

import (
	"errors"

	userdomain "github.com/codepnw/blog-api/internal/domains/user"
	"github.com/codepnw/blog-api/internal/handlers"
	"github.com/codepnw/blog-api/internal/middleware"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
	"github.com/codepnw/blog-api/internal/utils/errs"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc userusecase.Usecase
}

func NewUserHandler(uc userusecase.Usecase) *handler {
	return &handler{uc: uc}
}

// Create User
// @Summary Create User
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body userhandler.UserCreateReq true "User data"
// @Success 201 {object} userdomain.User
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users [post]
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

// Get User ID
// @Summary Get User ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} userdomain.User
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users/{user_id} [get]
func (h *handler) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyUserID)

	result, err := h.uc.GetUser(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return handlers.NotFound(ctx, err.Error())
		}
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Get Profile User
// @Summary Get Profile User
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userdomain.User
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users/me [get]
func (h *handler) GetProfile(ctx *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		return handlers.Unauthorized(ctx, err.Error())
	}

	result, err := h.uc.GetUser(ctx.Context(), user.ID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return handlers.NotFound(ctx, err.Error())
		}
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Get All User
// @Summary Get All User
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} []userdomain.User
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users [get]
func (h *handler) GetAllUsers(ctx *fiber.Ctx) error {
	result, err := h.uc.GetAllUsers(ctx.Context())
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Update User
// @Summary Update User
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Param data body userhandler.UserUpdateReq true "User data"
// @Success 200 {object} userdomain.User
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users/{user_id} [patch]
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
		if errors.Is(err, errs.ErrUserNotFound) {
			return handlers.NotFound(ctx, err.Error())
		}
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, result)
}

// Delete User
// @Summary Delete User
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} userdomain.User
// @Failure 404 {object} handlers.NotFoundRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /users/{user_id} [delete]
func (h *handler) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params(handlers.ParamKeyUserID)

	if err := h.uc.DeleteUser(ctx.Context(), id); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return handlers.NotFound(ctx, err.Error())
		}
		return handlers.InternalServerError(ctx, err)
	}
	return handlers.NoContent(ctx)
}

// ---- Auth ------

// Auth Register
// @Summary Auth Register
// @Tags auth
// @Accept json
// @Produce json
// @Param data body userhandler.UserCreateReq true "Register data"
// @Success 201 {object} userusecase.AuthResponse
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /auth/register [post]
func (h *handler) Register(ctx *fiber.Ctx) error {
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
	response, err := h.uc.Register(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Created(ctx, response)
}

// Auth Login
// @Summary Auth Login
// @Tags auth
// @Accept json
// @Produce json
// @Param data body userhandler.UserLoginReq true "Login data"
// @Success 200 {object} userusecase.AuthResponse
// @Failure 400 {object} handlers.BadRequestRes
// @Failure 500 {object} handlers.InternalServerErrRes
// @Router /auth/login [post]
func (h *handler) Login(ctx *fiber.Ctx) error {
	req := new(UserLoginReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return handlers.BadRequest(ctx, err.Error())
	}

	input := &userdomain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}
	response, err := h.uc.Login(ctx.Context(), input)
	if err != nil {
		return handlers.InternalServerError(ctx, err)
	}

	return handlers.Success(ctx, response)
}

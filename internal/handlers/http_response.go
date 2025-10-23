package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// -------------- Swagger Response --------------

type EmptyRes struct{}

type BadRequestRes struct {
	Message string `json:"message"`
}

type UnauthorizedRes struct {
	Message string `json:"message"`
}

type NotFoundRes struct {
	Message string `json:"message"`
}

type ForbiddenRes struct {
	Message string `json:"message"`
}

type InternalServerErrRes struct {
	Error string `json:"error"`
}

// -------------- JSON Response --------------

func Success(ctx *fiber.Ctx, data any) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"data": data,
	})
}

func Created(ctx *fiber.Ctx, data any) error {
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"data": data,
	})
}

func NoContent(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNoContent).JSON(nil)
}

func BadRequest(ctx *fiber.Ctx, message string) error {
	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message": message,
	})
}

func Unauthorized(ctx *fiber.Ctx, message string) error {
	return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
		"message": message,
	})
}

func NotFound(ctx *fiber.Ctx, message string) error {
	return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
		"message": message,
	})
}

func Forbidden(ctx *fiber.Ctx, message string) error {
	return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{
		"message": message,
	})
}

func InternalServerError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
		"error": err.Error(),
	})
}

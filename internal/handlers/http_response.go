package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

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

func InternalServerError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
		"error": err,
	})
}

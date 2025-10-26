package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func NewErrorResponse(ctx *fiber.Ctx, code int, errType, message string, details any) error {
	return ctx.Status(code).JSON(&fiber.Map{
		"error": ErrorResponse{
			Code:    code,
			Type:    errType,
			Message: message,
			Details: details,
		},
	})
}

func BadRequest(ctx *fiber.Ctx, message string) error {
	return NewErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func Unauthorized(ctx *fiber.Ctx, message string) error {
	return NewErrorResponse(ctx, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func NotFound(ctx *fiber.Ctx, message string) error {
	return NewErrorResponse(ctx, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func Forbidden(ctx *fiber.Ctx, message string) error {
	return NewErrorResponse(ctx, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func InternalServerError(ctx *fiber.Ctx, err error) error {
	return NewErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "", err.Error())
}

// ---------- Success Response -----------

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewSuccessResponse(ctx *fiber.Ctx, code int, message string, data any) error {
	return ctx.Status(code).JSON(&fiber.Map{
		"response": SuccessResponse{
			Code:    code,
			Message: message,
			Data:    data,
		},
	})
}

func Success(ctx *fiber.Ctx, data any) error {
	return NewSuccessResponse(ctx, http.StatusOK, "success", data)
}

func Created(ctx *fiber.Ctx, data any) error {
	return NewSuccessResponse(ctx, http.StatusCreated, "created successfully", data)
}

func NoContent(ctx *fiber.Ctx) error {
	return NewSuccessResponse(ctx, http.StatusNoContent, "no content", nil)
}

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

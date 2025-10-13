package routes

import (
	"database/sql"
	"errors"

	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	Prefix string     `validate:"required"`
	APP    *fiber.App `validate:"required"`
	DB     *sql.DB    `validate:"required"`
}

func RegisterRoutes(cfg *RouteConfig) (*RouteConfig, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, errors.New("required: prefix, app, db")
	}
	return cfg, nil
}

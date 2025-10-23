package routes

import (
	"database/sql"
	"errors"

	"github.com/codepnw/blog-api/internal/handlers/docs"
	_ "github.com/codepnw/blog-api/internal/handlers/post"
	"github.com/codepnw/blog-api/internal/middleware"
	jwttoken "github.com/codepnw/blog-api/internal/utils/jwt"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

type RouteConfig struct {
	Prefix string                    `validate:"required"`
	APP    *fiber.App                `validate:"required"`
	DB     *sql.DB                   `validate:"required"`
	Token  *jwttoken.JWTToken        `validate:"required"`
	Mid    *middleware.AppMiddleware `validate:"required"`
}

func RegisterRoutes(cfg *RouteConfig) (*RouteConfig, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, errors.New("required: prefix, app, db, token")
	}

	// Save Token in Local storage
	// Disable if production mode
	saveToken := true

	docs.SwaggerInfo.BasePath = cfg.Prefix
	cfg.APP.Get(cfg.Prefix+"/swagger/*", fiberSwagger.New(fiberSwagger.Config{
		PersistAuthorization: saveToken,
	}))

	return cfg, nil
}

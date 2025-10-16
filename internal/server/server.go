package server

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/config"
	"github.com/codepnw/blog-api/internal/database"
	"github.com/codepnw/blog-api/internal/middleware"
	"github.com/codepnw/blog-api/internal/server/routes"
	jwttoken "github.com/codepnw/blog-api/internal/utils/jwt"
	"github.com/gofiber/fiber/v2"
)

func Run(envPath string) error {
	// Load Config
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		return err
	}

	// Connect Database
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("connect database failed: %w", err)
	}
	defer db.Close()

	// Init JWT Token
	token, err := jwttoken.InitJWT(cfg)
	if err != nil {
		return err
	}

	// Init Middleware
	mid, err := middleware.InitMiddleware(token)
	if err != nil {
		return err
	}

	app := fiber.New()

	// Register Routes
	routesConfig := &routes.RouteConfig{
		Prefix: fmt.Sprintf("/api/v%d", cfg.APP.Version),
		APP:    app,
		DB:     db,
		Token:  token,
		Mid:    mid,
	}
	r, err := routes.RegisterRoutes(routesConfig)
	if err != nil {
		return err
	}
	// Init Routes
	r.CategoryRoutes()
	r.PostRoutes()
	r.UserRoutes()

	return app.Listen(fmt.Sprintf(":%d", cfg.APP.Port))
}

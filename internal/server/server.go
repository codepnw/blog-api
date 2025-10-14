package server

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/config"
	"github.com/codepnw/blog-api/internal/database"
	"github.com/codepnw/blog-api/internal/server/routes"
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

	app := fiber.New()

	// Register Routes
	routesConfig := &routes.RouteConfig{
		Prefix: fmt.Sprintf("/api/v%d", cfg.APP.Version),
		APP:    app,
		DB:     db,
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

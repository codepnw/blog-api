package server

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/config"
	"github.com/codepnw/blog-api/internal/database"
	"github.com/codepnw/blog-api/internal/middleware"
	"github.com/codepnw/blog-api/internal/server/routes"
	jwttoken "github.com/codepnw/blog-api/internal/utils/jwt"
	"github.com/codepnw/blog-api/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Run(envPath string) error {
	// Load Config
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		return err
	}

	// Logger
	logger.Init(cfg.APP.Mode)

	// Connect Database
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		logger.Error("server.Run: connect database", "error", err)
		return fmt.Errorf("connect database failed")
	}
	defer db.Close()

	// Init JWT Token
	token, err := jwttoken.InitJWT(cfg)
	if err != nil {
		logger.Error("server.Run: jwt init", "error", err)
		return err
	}

	// Init Middleware
	mid, err := middleware.InitMiddleware(token)
	if err != nil {
		logger.Error("server.Run: middleware init", "error", err)
		return err
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	// Register Routes
	routesConfig := &routes.RouteConfig{
		Mode:   cfg.APP.Mode,
		Prefix: fmt.Sprintf("/api/v%d", cfg.APP.Version),
		APP:    app,
		DB:     db,
		Token:  token,
		Mid:    mid,
	}
	r, err := routes.RegisterRoutes(routesConfig)
	if err != nil {
		logger.Error("server.Run: register routes", "error", err)
		return err
	}
	// Init Routes
	r.CategoryRoutes()
	r.PostRoutes()
	r.UserRoutes()
	r.CommentRoutes()

	port := fmt.Sprintf(":%d", cfg.APP.Port)
	url := fmt.Sprintf("%s%s%s", cfg.APP.Host, port, routesConfig.Prefix)
	logger.Info(fmt.Sprintf("server running at %s", url))
	logger.Info(fmt.Sprintf("api docs at %s/docs", url))

	if err := app.Listen(port); err != nil {
		logger.Error("server.Run: app listen", "error", err)
		return err
	}
	return nil
}

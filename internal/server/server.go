package server

import (
	"fmt"
	"log"

	"github.com/codepnw/blog-api/internal/config"
	"github.com/codepnw/blog-api/internal/database"
	"github.com/gofiber/fiber/v2"
)

func Run(envPath string) error {
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		return err
	}

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("connect database failed: %w", err)
	}
	defer db.Close()

	app := fiber.New()

	port := fmt.Sprintf(":%d", cfg.APP.Port)
	if err = app.Listen(port); err != nil {
		return err
	}

	log.Printf("server running at port %s", port)
	return nil
}

package database

import (
	"database/sql"
	"fmt"

	"github.com/codepnw/blog-api/internal/config"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *config.EnvConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database failed: %w", err)
	}

	return db, db.Ping()
}

package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/blog-api/internal/utils/validate"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	APP APPConfig `envPrefix:"APP_"`
	DB  DBConfig  `envPrefix:"DB_"`
}

type APPConfig struct {
	Port    int `env:"PORT" envDefault:"4000"`
	Version int `env:"VERSION" envDefault:"1"`
}

type DBConfig struct {
	User     string `env:"USER" validate:"required"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"NAME" validate:"required"`
	Host     string `env:"HOST" envDefault:"127.0.0.1" validate:"required"`
	Port     int    `env:"PORT" envDefault:"5432" validate:"required"`
	SSLMode  string `env:"SSLMODE" envDefault:"full-verify"`
}

func LoadConfig(path string) (*EnvConfig, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("load env failed: %w", err)
	}

	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate env failed: %w", err)
	}

	return cfg, nil
}

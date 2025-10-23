package main

import (
	"log"

	"github.com/codepnw/blog-api/internal/server"
)

const envPath = "dev.env"

// @title Blog API
// @version 1.0
// @description Clean Architecture Blog API
// @host localhost:4000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := server.Run(envPath); err != nil {
		log.Fatal(err)
	}
}

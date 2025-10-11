package main

import (
	"log"

	"github.com/codepnw/blog-api/internal/server"
)

const envPath = "dev.env"

func main() {
	if err := server.Run(envPath); err != nil {
		log.Fatal(err)
	}
}

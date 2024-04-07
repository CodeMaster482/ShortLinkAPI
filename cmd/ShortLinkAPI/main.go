package main

import (
	"fmt"
	"log"

	"ShortLinkAPI/config"
	"ShortLinkAPI/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}

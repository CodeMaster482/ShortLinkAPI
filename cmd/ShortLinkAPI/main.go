package main

import (
	"log"

	"github.com/CodeMaster482/ShortLinkAPI/config"
	"github.com/CodeMaster482/ShortLinkAPI/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Config error: %s", err)
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}

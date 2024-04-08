package main

import (
	"flag"
	"log"

	"github.com/CodeMaster482/ShortLinkAPI/config"
	"github.com/CodeMaster482/ShortLinkAPI/internal/app"

	"github.com/joho/godotenv"
)

var useRedis bool

func init() {
	flag.BoolVar(&useRedis, "in-memo", false, "Set to true if use Redis, otherwise PostgreSQL.")
	flag.Parse()
}

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

	cfg.UseRedis = useRedis

	app.Run(cfg)
}

package app

import (
	"crypto"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ShortLinkAPI/config"
	"ShortLinkAPI/pkg/generator"
	"ShortLinkAPI/pkg/httpserver"
	"ShortLinkAPI/pkg/logger"
	"ShortLinkAPI/pkg/postgres"

	linkHandler "ShortLinkAPI/internal/delivery/http/handler"
	linkRepository "ShortLinkAPI/internal/repository/postgres"
	linkUsecase "ShortLinkAPI/internal/usecase"

	"github.com/gin-gonic/gin"
)

// @title Go ShortLinkAPI
// @version 1.0
// @description Golang REST API for creating, handling short links.
// @contact.name Grigory Kovalenko
// @contact.url https://github.com/CodeMaster482
// @contact.email grigorikovalenko@gmail.com
// @BasePath /
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(
		cfg.PG.Host,
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Name,
		cfg.PG.Port,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	generator := generator.NewGenerator(
		generator.WithAlphabet(cfg.LinkGen.Alphabet),
		generator.WithHashFunc(crypto.MD5),
		generator.WithLength(cfg.LinkGen.Length),
	)

	// Use case
	lr := linkRepository.NewLinkStorage(pg.Pool)
	lu := linkUsecase.NewLinkService(cfg, lr, generator)
	lh := linkHandler.NewLinkHandler(lu)

	// HTTP Server
	r := gin.New()
	api := r.Group("/api/v1")

	api.POST("/url", lh.CreateLink)
	api.GET("/url/:key", lh.GetLink)

	httpServer := httpserver.New(
		r,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.ReadTimeout*time.Second),
		httpserver.WriteTimeout(cfg.HTTP.WriteTimeout*time.Second),
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

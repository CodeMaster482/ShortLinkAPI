package app

import (
	"context"
	"crypto"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	linkGrpcHandler "github.com/CodeMaster482/ShortLinkAPI/internal/delivery/grpc"
	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/grpc/generated"
	linkHandler "github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/handler"
	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/middleware"
	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	linkSQLRepo "github.com/CodeMaster482/ShortLinkAPI/internal/repository/postgres"
	linkRedisRepo "github.com/CodeMaster482/ShortLinkAPI/internal/repository/redis"
	linkUsecase "github.com/CodeMaster482/ShortLinkAPI/internal/usecase"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"github.com/CodeMaster482/ShortLinkAPI/config"
	"github.com/CodeMaster482/ShortLinkAPI/pkg/generator"
	"github.com/CodeMaster482/ShortLinkAPI/pkg/httpserver"
	"github.com/CodeMaster482/ShortLinkAPI/pkg/logger"
	"github.com/CodeMaster482/ShortLinkAPI/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type LinkRepository interface {
	GetLink(ctx context.Context, token string) (*model.Link, error)
	StoreLink(ctx context.Context, link *model.Link) error
	StartRecalculation(interval time.Duration, deleted chan []string)
}

func addPingRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/ping")

	ping.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}

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
	var lr LinkRepository

	if cfg.UseRedis {
		cli := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		})

		_, err := cli.Ping(context.Background()).Result()
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - redis.ping: %w", err))
		}

		lr = linkRedisRepo.NewLinkStorage(cli)
	} else {
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

		lr = linkSQLRepo.NewLinkStorage(pg.Pool)
	}

	g := generator.NewGenerator(
		generator.WithAlphabet(cfg.LinkGen.Alphabet),
		generator.WithHashFunc(crypto.MD5),
		generator.WithLength(cfg.LinkGen.Length),
	)

	// Use case
	lu := linkUsecase.NewLinkService(cfg, lr, g)
	lh := linkHandler.NewLinkHandler(lu)

	// HTTP Server
	r := gin.New()
	base := r.Group("/")
	addPingRoutes(base)
	api := r.Group("/api/v1")

	api.Use(middleware.ErrorMiddleware())
	api.Use(middleware.RequestTimeout(500 * time.Millisecond))
	api.Use(gin.Logger(), gin.Recovery())

	api.POST("/url", lh.CreateLink)
	api.GET("/url/:key", lh.GetLink)

	httpServer := httpserver.New(
		r,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.ReadTimeout),
		httpserver.WriteTimeout(cfg.HTTP.WriteTimeout),
	)

	grpcHandler := linkGrpcHandler.NewLinkHandler(lu)
	grpcServer := grpc.NewServer()
	generated.RegisterShortLinkServiceServer(grpcServer, grpcHandler)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		l.Fatal(err)
	}
	defer grpcListener.Close()

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			l.Fatal(err)
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	grpcServer.GracefulStop()

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

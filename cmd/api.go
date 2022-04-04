package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/message"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/service"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/logkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/migrationkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/pgkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/rediskit"
	"go.uber.org/zap"
)

type APIArgs struct {
	rediskit.Redisconfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	pgkit.PGConfig       `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
	logkit.LoggerConfig  `group:"logger" namespace:"logger" env-namespace:"LOGGER"`
}

const gracefulWaitSecond = 5

func runAPI() {
	ctx := context.Background()

	var args APIArgs
	if _, err := flags.NewParser(&args, flags.Default).Parse(); err != nil {
		log.Fatal("failed to parse flag", err.Error())
	}

	logger := logkit.NewLogger(&args.LoggerConfig)
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal("fail to sync logger", err.Error())
		}
	}()

	ctx = logger.WithContext(ctx)

	redisClient := rediskit.NewRedisClient(ctx, &args.Redisconfig)
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Fatal("failed to close redis client", zap.Error(err))
		}
	}()

	pgClient := pgkit.NewPGClient(ctx, &args.PGConfig)
	defer func() {
		if err := pgClient.Close(); err != nil {
			logger.Fatal("failed to close postgres client", zap.Error(err))
		}
	}()

	migrationConf := &migrationkit.MigrationConfig{
		Source: "file://migrations",
		URL:    args.PGConfig.URL,
	}

	migration := migrationkit.NewMigration(ctx, migrationConf)
	defer func() {
		if err := migration.Close(); err != nil {
			logger.Fatal("failed to close migration", zap.Error(err))
		}
	}()

	if err := migration.Up(); err != nil {
		logger.Fatal("failed to call migration up", zap.Error(err))
	}

	pgShortenerDAO := dao.NewPGShortenerDAO(pgClient)
	shortenerDAO := dao.NewRedisShortenerDAO(redisClient, pgShortenerDAO)
	svc := service.NewService(shortenerDAO)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(message.SuccessRequest, gin.H{
			"status": "ok",
		})
	})

	router.GET("/get", func(c *gin.Context) {
		svc.GetURL(c)
	})

	router.POST("/shorten", func(c *gin.Context) {
		svc.ShortenURL(c)
	})

	logger.Info("listening to port 8080")

	srv := &http.Server{
		Addr:    ":8008",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal("fail to serve the request", zap.Error(err))
		}
	}()

	// gracefully shutdown with timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("get shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulWaitSecond*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("fail to shutdown server", zap.Error(err))
	}
	logger.Info("server has already shutdown")
}

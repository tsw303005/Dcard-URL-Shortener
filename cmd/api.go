package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/service"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/logkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/pgkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/rediskit"
	"go.uber.org/zap"
)

type APIArgs struct {
	rediskit.Redisconfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	pgkit.PGConfig       `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
	logkit.LoggerConfig  `group:"logger" namespace:"logger" env-namespace:"LOGGER"`
}

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

	pgShortenerDAO := dao.NewPGShortenerDAO(pgClient)
	shortenerDAO := dao.NewRedisShortenerDAO(redisClient, pgShortenerDAO)
	svc := service.NewService(shortenerDAO)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.GET("/get", func(c *gin.Context) {
		svc.GetUrl(c)
	})

	r.POST("/shorten", func(c *gin.Context) {
		svc.ShortenUrl(c)
	})

	logger.Info("listening to port 8080")
	err := r.Run(":8080")

	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/service"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/pgkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/rediskit"
)

type APIArgs struct {
	rediskit.Redisconfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	pgkit.PGConfig       `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
}

func runAPI() {
	ctx := context.Background()

	var args APIArgs
	if _, err := flags.NewParser(&args, flags.Default).Parse(); err != nil {
		log.Fatal("failed to parse flag", err.Error())
	}

	redisClient := rediskit.NewRedisClient(ctx, &args.Redisconfig)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatal("failed to close redis client", err)
		}
	}()

	pgClient := pgkit.NewPGClient(ctx, &args.PGConfig)
	defer func() {
		if err := pgClient.Close(); err != nil {
			log.Fatal("failed to close postgres client", err)
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

	err := r.Run(":8080")

	if err != nil {
		log.Fatal(err)
	}
}

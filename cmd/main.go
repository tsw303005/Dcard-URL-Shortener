package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/service"
)

func main() {
	runAPI()
}

func runAPI() {
	ctx := context.Background()

	testDAO := dao.NewTestDAO()
	svc := service.NewService(testDAO)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
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

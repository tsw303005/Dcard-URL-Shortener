package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/message"
)

type Service struct {
	URLDAO dao.ShortenerDAO
}

func NewService(URLDAO dao.ShortenerDAO) *Service {
	return &Service{
		URLDAO: URLDAO,
	}
}

func (s *Service) GetUrl(c *gin.Context) {
	shorten_url := c.Query("shorten_url")

	shortener, err := s.URLDAO.Get(c.Request.Context(), &dao.Shortener{
		ShortenUrl: shorten_url,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"request": "get url",
		})
		log.Fatal(err)
	}

	c.Redirect(302, shortener.Url)
}

func (s *Service) ShortenUrl(c *gin.Context) {
	var req message.ShortenUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"request": "shorten url",
		})
	}

	url_id, shorten_url, err := s.URLDAO.Shorten(c.Request.Context(), &dao.Shortener{
		Url:       req.Url,
		ExpiredAt: req.ExpiredAt,
	})

	if err != nil {
		log.Fatal(err)
	}

	resp := message.ShortenUrlResponse{
		Id:         url_id,
		ShortenUrl: shorten_url,
	}

	c.JSON(200, resp)
}

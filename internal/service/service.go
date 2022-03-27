package service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/message"
)

type Service struct {
	urlDAO dao.ShortenerDAO
}

func NewService(urlDAO dao.ShortenerDAO) *Service {
	return &Service{
		urlDAO: urlDAO,
	}
}

func (s *Service) GetURL(c *gin.Context) {
	shortenURL := c.Query("shorten_url")

	shortener, err := s.urlDAO.Get(c.Request.Context(), &dao.Shortener{
		ShortenURL: shortenURL,
	})

	if err == dao.ErrExpiredat {
		c.JSON(message.URLExpired, gin.H{
			"error":   "this shorten url has already expired",
			"request": "redirct url request",
		})
		return
	} else if err != nil {
		c.JSON(message.ServiceError, gin.H{
			"error":   "internal server error",
			"request": "redirct url request",
		})
		log.Fatal("internal service error")
	}

	c.Redirect(message.URLRedirect, shortener.URL)
}

func (s *Service) ShortenURL(c *gin.Context) {
	var req message.ShortenURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(message.BadRequest, gin.H{
			"error":   "argument error",
			"request": "shorten url request",
		})
		return
	}

	urlID, shortenURL, err := s.urlDAO.Shorten(c.Request.Context(), &dao.Shortener{
		URL:       req.URL,
		ExpiredAt: req.ExpiredAt,
	})

	if err != nil {
		c.JSON(message.ServiceError, gin.H{
			"error":   "internal server error",
			"request": "shorten url request",
		})
		return
	}

	resp := message.ShortenURLResponse{
		ID:         urlID,
		ShortenURL: shortenURL,
	}

	c.JSON(message.SuccessRequest, resp)
}

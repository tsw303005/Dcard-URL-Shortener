package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	urlID, err := uuid.Parse(c.Param("url_id"))

	if err != nil || urlID == uuid.Nil {
		c.JSON(message.BadRequest, gin.H{
			"error":   "this url id has a wrong format",
			"request": "redirect url request",
		})
		return
	}

	shortener, err := s.urlDAO.Get(c.Request.Context(), &dao.Shortener{
		ID: urlID,
	})

	if err == dao.ErrExpiredat {
		c.JSON(message.URLExpired, gin.H{
			"error":   "this shorten url has already expired",
			"request": "redirect url request",
		})
		return
	} else if err == dao.ErrShortenURLNotFound {
		c.JSON(message.BadRequest, gin.H{
			"error":   "this shorten url not found",
			"request": "redirect url request",
		})
		return
	} else if err != nil {
		c.JSON(message.ServiceError, gin.H{
			"error":   "internal server error",
			"request": "redirect url request",
		})
		return
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
		ShortenURL: c.Request.Host,
		URL:        req.URL,
		ExpiredAt:  req.ExpiredAt,
	})

	if err == dao.ErrShortenURLFail {
		c.JSON(message.BadRequest, gin.H{
			"error":   "expiredAt has already expired",
			"request": "shorten url request",
		})
		return
	} else if err != nil {
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

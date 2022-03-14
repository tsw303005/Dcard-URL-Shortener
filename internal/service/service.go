package service

import (
	"context"

	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
)

type service struct {
	URLDAO dao.URL
}

func NewService(URLDAO dao.URL) *service {
	return &service{
		URLDAO: URLDAO,
	}
}

func (s *service) ShortenURL(ctx context.Context)

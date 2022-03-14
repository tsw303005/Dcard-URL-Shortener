package dao

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID         uuid.UUID
	ExpiredAt  time.Time
	URL        string
	ShortenURL string
}

type URLDAO interface {
	ShortenURL(ctx context.Context, url string, expiredAt time.Time) (URL, error)
	GetURLByID(ctx context.Context, ID uuid.UUID) (URL, error)
}

var (
	ExpiredURLError = errors.New("id has alread expired")
)

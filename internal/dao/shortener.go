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
	Shorten(ctx context.Context, url string, expiredAt time.Time) (string, error)
	Get(ctx context.Context, ID uuid.UUID) (URL, error)
}

var (
	ExpiredURLError = errors.New("id has alread expired")
)

package dao

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID         uuid.UUID
	URL        string
	ShortenURL string
	ExpiredAt  time.Time
}

type URLDAO interface {
	Shorten(ctx context.Context, url string, expiredAt time.Time) (string, string, error)
	Get(ctx context.Context, ID uuid.UUID) (string, error)
}

var (
	ExpiredURLError = errors.New("id has alread expired")
)

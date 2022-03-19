package dao

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Shortener struct {
	ID         uuid.UUID
	Url        string
	ShortenUrl string
	ExpiredAt  string
}

type ShortenerDAO interface {
	Shorten(ctx context.Context, req *Shortener) (uuid.UUID, string, error)
	Get(ctx context.Context, req *Shortener) (*Shortener, error)
}

var (
	ErrExpiredat          = errors.New("id has already expired")
	ErrShortenURLNotFound = errors.New("shorten url not found")
	ErrShortenURLFail     = errors.New("fail to shorten url")
)

func getShortenerUrl(shortenerUrl string) string {
	return "getUrl:" + shortenerUrl
}

package dao

import (
	"errors"
	"time"
)

type URL struct {
	ID         string
	Url        string
	ShortenUrl string
	ExpiredAt  time.Time
}

type URLDAO interface {
	Shorten(req URL) (string, string, error)
	Get(req URL) (string, error)
}

var (
	ExpiredURLError = errors.New("id has alread expired")
)

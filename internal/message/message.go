package message

import (
	"time"
)

type ShortenUrlRequest struct {
	Url       string
	ExpiredAt time.Time
}

type ShortenUrlResponse struct {
	Id         string
	ShortenUrl string
}

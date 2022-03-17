package message

import (
	"time"
)

type ShortenUrlRequest struct {
	Url       string    `json:"url"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type ShortenUrlResponse struct {
	Id         string `json:"id"`
	ShortenUrl string `json:"shortUrl"`
}

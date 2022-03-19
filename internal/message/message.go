package message

import (
	"github.com/google/uuid"
)

type ShortenUrlRequest struct {
	Url       string `json:"url"`
	ExpiredAt string `json:"expiredAt"`
}

type ShortenUrlResponse struct {
	Id         uuid.UUID `json:"id"`
	ShortenUrl string    `json:"shortUrl"`
}

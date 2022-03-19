package message

import (
	"time"

	"github.com/google/uuid"
)

type ShortenUrlRequest struct {
	Url       string    `json:"url"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type ShortenUrlResponse struct {
	Id         uuid.UUID `json:"id"`
	ShortenUrl string    `json:"shortUrl"`
}

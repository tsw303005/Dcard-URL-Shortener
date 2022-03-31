package message

import (
	"github.com/google/uuid"
)

type ShortenURLRequest struct {
	URL       string `json:"url"`
	ExpiredAt string `json:"expiredAt"`
}

type ShortenURLResponse struct {
	ID         uuid.UUID `json:"id"`
	ShortenURL string    `json:"shortenURL"`
}

package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
)

type service struct {
	URLDAO dao.URLDAO
}

type shortenResponse struct {
	code     int
	id       string
	shortUrl string
}

func NewService(URLDAO dao.URLDAO) *service {
	return &service{
		URLDAO: URLDAO,
	}
}

func (s *service) ShortenURL(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req dao.URL

	// parse request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Fatal(err)
	}

	// shorten url
	url, err := s.URLDAO.Shorten(ctx, req.URL, req.ExpiredAt)

	if err != nil {
		log.Fatal(err)
	}

	res := shortenResponse{
		code:     200,
		id:       url,
		shortUrl: "http://localhost/" + url,
	}

	json.NewEncoder(w).Encode(res)
}

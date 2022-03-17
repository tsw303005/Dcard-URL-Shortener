package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
)

type Service struct {
	URLDAO dao.URLDAO
}

type shortenResponse struct {
	code     int
	id       string
	shortUrl string
}

type getResponse struct {
	code        int
	originalUrl string
}

func NewService(URLDAO dao.URLDAO) *Service {
	return &Service{
		URLDAO: URLDAO,
	}
}

func (s *Service) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req dao.URL

	// parse request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Fatal(err)
	}

	// shorten url
	id, url, err := s.URLDAO.Shorten(r.Context(), req.URL, req.ExpiredAt)

	if err != nil {
		log.Fatal(err)
	}

	res := shortenResponse{
		code:     200,
		id:       id,
		shortUrl: url,
	}

	json.NewEncoder(w).Encode(res)
}

func (s *Service) GetURL(w http.ResponseWriter, r *http.Request) {
	var req dao.URL

	// parse request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Fatal(err)
	}

	url, err := s.URLDAO.Get(r.Context(), req.ID)

	if err != nil {
		log.Fatal(err)
	}

	res := getResponse{
		code:        200,
		originalUrl: url,
	}

	json.NewEncoder(w).Encode(res)
}

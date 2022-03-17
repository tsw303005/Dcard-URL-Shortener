package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/message"
)

type Service struct {
	URLDAO dao.URLDAO
}

func NewService(URLDAO dao.URLDAO) *Service {
	return &Service{
		URLDAO: URLDAO,
	}
}

func (s *Service) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req message.ShortenUrlRequest

	// parse request
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		log.Fatal(err)
	}

	// shorten url
	url_id, shorten_url, err := s.URLDAO.Shorten(r.Context(), dao.URL{
		Url:       req.Url,
		ExpiredAt: req.ExpiredAt,
	})

	if err != nil {
		log.Fatal(err)
	}

	res := message.ShortenUrlResponse{
		Id:         url_id,
		ShortenUrl: shorten_url,
	}

	js, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *Service) GetURL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("ID")

	url, err := s.URLDAO.Get(r.Context(), dao.URL{
		ID: id,
	})

	if err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(url)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

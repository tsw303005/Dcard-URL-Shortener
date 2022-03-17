package main

import (
	"log"
	"net/http"

	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/service"
)

func main() {
	runAPI()
}

func runAPI() {
	var URLDAO dao.URLDAO

	svc := service.NewService(URLDAO)

	http.HandleFunc("/get", svc.GetURL)
	http.HandleFunc("/shorten", svc.ShortenURL)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"net/http"

	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/service"
)

func main() {

}

func runAPI() {

	var URLDAO dao.URLDAO

	svc := service.NewService(URLDAO)

	http.HandleFunc("/get", svc.GetURL)

}

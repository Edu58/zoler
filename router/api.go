package router

import (
	"log"
	"net/http"

	"github.com/Edu58/zoler/controller"
)

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/crawl", controller.Crawler)

	log.Fatal(http.ListenAndServe(":4500", mux))
}

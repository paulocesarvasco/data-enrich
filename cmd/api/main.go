package main

import (
	"log"
	"net/http"

	cte "data-enrich/internal/constants"
	"data-enrich/internal/engine"

	"github.com/gorilla/mux"
)

func main() {
	enrich := engine.New()
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/input", enrich.Enrich()).Methods(http.MethodPost).Schemes("http").Headers("Content-Type", "application/json")
	r.HandleFunc("/get", enrich.Search()).Methods(http.MethodGet).Schemes("http")

	log.Println(cte.StartingServer)
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"log"
	"net/http"

	cte "data-enrich/internal/constants"
	"data-enrich/internal/database"
	"data-enrich/internal/engine"
	"data-enrich/internal/utils"

	"github.com/gorilla/mux"
)

func main() {

	var err error

	err = database.CreateDatabaseInstance("mongodb://mongodb:27017")
	if err != nil {
		log.Fatal(utils.WrapError(err, cte.ErrorToCreateDatabaseInstace))
	}

	err = database.CreateDatabaseCollection("CloudtrailRecords", "records")
	if err != nil {
		log.Fatal(utils.WrapError(err, cte.ErrorToCreateDatabaseCollection))
	}

	enrich := engine.New()
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/input", enrich.Enrich()).Methods(http.MethodPost).Schemes("http").Headers("Content-Type", "application/json")
	r.HandleFunc("/get", enrich.Search()).Methods(http.MethodGet).Schemes("http")

	log.Println(cte.StartingServer)
	log.Fatal(http.ListenAndServe(":8080", r))
}

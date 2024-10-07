package main

import (
	"log"

	"data-enrich/connection"
	cte "data-enrich/constants"
	"data-enrich/database"
	"data-enrich/utils"
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

	log.Println(cte.ConnectionWithDbEstablish)
	connection.HandleRequests()
}

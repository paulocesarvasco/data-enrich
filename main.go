package main

import (
	"fmt"

	"meli/connection"
	"meli/database"
)

func main() {

	var err error

	err = database.CreateDatabaseInstance("mongodb://localhost:27017")
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	err = database.CreateDatabaseCollection("CloudtrailRecords", "records")
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	connection.HandleRequests()
}

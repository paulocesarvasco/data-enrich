package main

import (
	"fmt"

	"meli/connection"
	cte "meli/constants"
	"meli/database"
	"meli/utils"
)

func main() {

	var err error

	err = database.CreateDatabaseInstance("mongodb://localhost:27017")
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrorToCreateDatabaseInstace))
	}

	err = database.CreateDatabaseCollection("CloudtrailRecords", "records")
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrorToCreateDatabaseCollection)) // output for debug
	}

	connection.HandleRequests()
}

package connection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	cte "meli/constants"
	"meli/database"
	"meli/models"
	"meli/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/input", saveDataOnDb).Methods(http.MethodPost).Schemes("http").Headers("Content-Type", "application/json")
	r.HandleFunc("/get", getData).Methods(http.MethodGet).Schemes("http")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func saveDataOnDb(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)

	var record models.CloudtrailData

	err := json.Unmarshal(reqBody, &record)
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrorToUnmarshallRequestBody)) // output for debug
	}

	country, err := utils.GetCountryFromIp(record.Records[0].SourceIPAddress)
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrorToRetriveCountryFromIp)) // output for debug
	}

	region, err := utils.GetCountryRegion(country)
	if err != nil {
		fmt.Printf("%+v\n", cte.ErrorToRetriveRegionName) // output for debug
	}

	enrichiment := models.Enrichment{
		Country: country,
		Region:  region,
	}

	record.Records[0].Enrichment = enrichiment

	err = database.SaveDataOnDatabase(record)
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrortoSaveDataOnDatabase)) // output for debug

	} else {
		w.Write([]byte("Ok\n"))
	}
}

func getData(w http.ResponseWriter, r *http.Request) {

	numRecords := 10 // hard coded for now
	records, err := database.RetriveLastDataFromDatabase(numRecords)
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrorToRetrieveRecordsFromDb)) // output for debug

	}

	byteRecords, err := json.Marshal(records)
	if err != nil {
		fmt.Printf("%+v\n", utils.WrapError(err, cte.ErrortoEncodeDataFromDatabase)) // output for debug

	} else {
		w.Write(byteRecords)
	}
}

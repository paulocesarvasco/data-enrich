package connection

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	cte "meli/constants"
	"meli/database"
	"meli/models"
	"meli/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// Set POST and GET endpoints and start rest server
func HandleRequests() {

	log.Println(cte.StartingServer)
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/input", saveDataOnDb).Methods(http.MethodPost).Schemes("http").Headers("Content-Type", "application/json")
	r.HandleFunc("/get", getData).Methods(http.MethodGet).Schemes("http")

	log.Fatal(http.ListenAndServe(":8080", r)) // Port hard-coded for now
}

// Parses http request retrieve IPSource, gets Country and Region, send new data to db
func saveDataOnDb(w http.ResponseWriter, r *http.Request) {

	// Read http request
	reqBody, _ := ioutil.ReadAll(r.Body)

	var record models.CloudtrailData

	// Convert request from json to go type
	err := json.Unmarshal(reqBody, &record)
	if err != nil {
		log.Println(utils.WrapError(err, cte.ErrorToUnmarshallRequestBody))
		return
	}

	// Get country name from IP
	country, err := utils.GetCountryFromIp(record.Records[0].SourceIPAddress)
	if err != nil {
		log.Fatal(utils.WrapError(err, cte.ErrorToRetriveCountryFromIp))
		return
	}

	// Get region name from country name
	region, err := utils.GetCountryRegion(country)
	if err != nil {
		log.Println(utils.WrapError(err, cte.ErrorToRetriveRegionName))
		return
	}

	enrichiment := models.Enrichment{
		Country: country,
		Region:  region,
	}

	// Insert new values into input data
	record.Records[0].Enrichment = enrichiment

	// Save changed input in the db
	err = database.SaveDataOnDatabase(record)
	if err != nil {
		log.Println(cte.ErrortoSaveDataOnDatabase)
		return
	}

	// Only to ensure data integrity, calculate hash value and log
	dataHash := sha256.Sum256(reqBody)
	id := base64.StdEncoding.EncodeToString(dataHash[:])

	logMessage := cte.DataSavedWithId + id
	log.Println(logMessage)
}

func getData(w http.ResponseWriter, r *http.Request) {

	numRecords := 10 // hard coded for now
	records, err := database.RetriveLastDataFromDatabase(numRecords)
	if err != nil {
		log.Println(err, cte.ErrorToRetrieveRecordsFromDb)
		return

	}

	// Convert data retrieved to json format
	byteRecords, err := json.Marshal(records)
	if err != nil {
		log.Println(cte.ErrortoEncodeDataFromDatabase)
		return
	} else {
		w.Write(byteRecords)
	}

	// Log the IP of the requester
	logMessage := cte.DataRetrieved + r.Host
	log.Println(logMessage)

}

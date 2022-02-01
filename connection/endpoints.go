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

func HandleRequests() {

	log.Println(cte.StartingServer)
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
		log.Println(utils.WrapError(err, cte.ErrorToUnmarshallRequestBody))
		return
	}

	country, err := utils.GetCountryFromIp(record.Records[0].SourceIPAddress)
	if err != nil {
		log.Fatal(utils.WrapError(err, cte.ErrorToRetriveCountryFromIp))
		return
	}

	region, err := utils.GetCountryRegion(country)
	if err != nil {
		log.Println(utils.WrapError(err, cte.ErrorToRetriveRegionName))
		return
	}

	enrichiment := models.Enrichment{
		Country: country,
		Region:  region,
	}

	record.Records[0].Enrichment = enrichiment

	err = database.SaveDataOnDatabase(record)
	if err != nil {
		log.Println(cte.ErrortoSaveDataOnDatabase)
		return
	}

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

	byteRecords, err := json.Marshal(records)
	if err != nil {
		log.Println(cte.ErrortoEncodeDataFromDatabase)
		return
	} else {
		w.Write(byteRecords)
	}

	logMessage := cte.DataRetrieved + r.RemoteAddr
	log.Println(logMessage)

}

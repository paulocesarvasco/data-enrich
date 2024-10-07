package engine

import (
	"crypto/sha256"
	cte "data-enrich/internal/constants"
	"data-enrich/internal/database"
	"data-enrich/internal/models"
	"data-enrich/internal/utils"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type API interface {
	Enrich() http.HandlerFunc
	Search() http.HandlerFunc
}

type api struct {
}

func New() API {
	return &api{}
}

// Enrich parses http request retrieve IPSource, gets Country and Region, send new data to db
func (a *api) Enrich() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read http request
		reqBody, _ := io.ReadAll(r.Body)

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
}

func (a *api) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
}

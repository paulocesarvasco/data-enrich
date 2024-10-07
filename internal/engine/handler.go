package engine

import (
	cte "data-enrich/internal/constants"
	"data-enrich/internal/database"
	"data-enrich/internal/models"
	"data-enrich/internal/utils"
	"encoding/json"
	"log"
	"net/http"
)

type API interface {
	Enrich() http.HandlerFunc
	Search() http.HandlerFunc
}

type api struct {
	db database.DB
}

func New() API {
	return &api{}
}

// Enrich parses http request retrieve IPSource, gets Country and Region, send new data to db
func (a *api) Enrich() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var record models.CloudtrailData
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			log.Print(utils.WrapError(err, cte.ErrorToUnmarshallRequestBody))
			http.Error(w, cte.ErrorToUnmarshallRequestBody, http.StatusBadRequest)
			return
		}

		if len(record.Records) == 0 {
			log.Print(cte.ErrorToUnmarshallRequestBody)
			http.Error(w, cte.ErrorMissedMandatoryFields, http.StatusBadRequest)
			return
		}

		// Get country name from IP
		country, err := utils.GetCountryFromIp(record.Records[0].SourceIPAddress)
		if err != nil {
			log.Print(utils.WrapError(err, cte.ErrorToRetriveCountryFromIp))
			http.Error(w, cte.ErrorToRetriveCountryFromIp, http.StatusInternalServerError)
			return
		}

		// Get region name from country name
		region, err := utils.GetCountryRegion(country)
		if err != nil {
			log.Print(utils.WrapError(err, cte.ErrorToRetriveRegionName))
			http.Error(w, cte.ErrorToRetriveRegionName, http.StatusInternalServerError)
			return
		}

		enrichiment := models.Enrichment{
			Country: country,
			Region:  region,
		}

		// Insert new values into input data
		record.Records[0].Enrichment = enrichiment

		// Save changed input in the db
		err = a.db.Save(r.Context(), record)
		if err != nil {
			log.Print(cte.ErrortoSaveDataOnDatabase)
			http.Error(w, cte.ErrortoSaveDataOnDatabase, http.StatusInternalServerError)
			return
		}
		log.Println("data saved.")
	}
}

func (a *api) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := a.db.RetriveLastregisters(r.Context(), cte.NUM_RECORDS)
		if err != nil {
			log.Print(err, cte.ErrorToRetrieveRecordsFromDb)
			http.Error(w, cte.ErrorToRetrieveRecordsFromDb, http.StatusInternalServerError)
			return
		}

		// Convert data retrieved to json format
		byteRecords, err := json.Marshal(records)
		if err != nil {
			log.Println(utils.WrapError(err, cte.ErrortoEncodeDataFromDatabase))
			http.Error(w, cte.ErrorToRetrieveRecordsFromDb, http.StatusInternalServerError)
			return
		} else {
			w.Write(byteRecords)
		}

		// Log the IP of the requester
		logMessage := cte.DataRetrieved + r.Host
		log.Println(logMessage)
	}
}

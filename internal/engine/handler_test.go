package engine

import (
	"bytes"
	"data-enrich/internal/constants"
	database "data-enrich/internal/database/mocks"
	enrich "data-enrich/internal/enrich/mocks"
	"data-enrich/internal/errors"
	"data-enrich/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestEnrich(t *testing.T) {
	tt := []struct {
		name              string
		requestBody       any
		saveDatabaseError error
		countryEnrich     struct {
			err     error
			country string
		}
		regionEnrich struct {
			err    error
			region string
		}
		expectedBody string
		expectedCode int
	}{
		{
			name:              "success",
			requestBody:       models.CloudtrailData{Records: []models.Records{{SourceIPAddress: "10.10.10.10"}}},
			saveDatabaseError: nil,
			countryEnrich: struct {
				err     error
				country string
			}{err: nil, country: "Brazil"},
			regionEnrich: struct {
				err    error
				region string
			}{err: nil, region: "South America"},
			expectedBody: "ACK",
			expectedCode: http.StatusOK,
		},
		{
			name:              "invalid body",
			requestBody:       "invalid data",
			saveDatabaseError: nil,
			countryEnrich: struct {
				err     error
				country string
			}{},
			regionEnrich: struct {
				err    error
				region string
			}{},
			expectedBody: constants.ErrorToUnmarshallRequestBody,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:              "request without records",
			requestBody:       models.CloudtrailData{},
			saveDatabaseError: nil,
			countryEnrich: struct {
				err     error
				country string
			}{},
			regionEnrich: struct {
				err    error
				region string
			}{},
			expectedBody: constants.ErrorMissedMandatoryFields,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:              "get country error",
			requestBody:       models.CloudtrailData{Records: []models.Records{{SourceIPAddress: "10.10.10.10"}}},
			saveDatabaseError: nil,
			countryEnrich: struct {
				err     error
				country string
			}{err: errors.ErrorToRetriveCountryFromIp},
			regionEnrich: struct {
				err    error
				region string
			}{},
			expectedBody: constants.ErrorToRetriveCountryFromIp,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:              "get region error",
			requestBody:       models.CloudtrailData{Records: []models.Records{{SourceIPAddress: "10.10.10.10"}}},
			saveDatabaseError: nil,
			countryEnrich: struct {
				err     error
				country string
			}{err: nil, country: "Brazil"},
			regionEnrich: struct {
				err    error
				region string
			}{err: errors.ErrorToRetriveRegionName},
			expectedBody: constants.ErrorToRetriveRegionName,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:              "save database fail",
			requestBody:       models.CloudtrailData{Records: []models.Records{{SourceIPAddress: "10.10.10.10"}}},
			saveDatabaseError: errors.ErrorDatabaseOperationSave,
			countryEnrich: struct {
				err     error
				country string
			}{err: nil, country: "Brazil"},
			regionEnrich: struct {
				err    error
				region string
			}{err: nil, region: "South America"},
			expectedBody: constants.ErrorDatabaseOperationSave,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		db := database.NewClient()
		db.SetSaveError(tc.saveDatabaseError)
		enrh := enrich.NewEnrichService()
		enrh.SetCountryResponse(tc.countryEnrich.country)
		enrh.SetCountrySearchError(tc.countryEnrich.err)
		enrh.SetRegionResponse(tc.regionEnrich.region)
		enrh.SetSearchRegionError(tc.regionEnrich.err)

		a := api{
			db:     db,
			enrich: enrh,
		}
		r := mux.NewRouter()
		r.Handle("/enrich", a.Enrich())

		ts := httptest.NewServer(r)
		defer ts.Close()

		rawBody, _ := json.Marshal(tc.requestBody)
		readerBody := bytes.NewReader(rawBody)
		resp, err := http.Post(ts.URL+"/enrich", "", readerBody)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != tc.expectedCode {
			t.Errorf("TEST [%s] FAILED: expected code [%d] but received [%d]", tc.name, tc.expectedCode, resp.StatusCode)
		}

	}
}

func TestSearch(t *testing.T) {
	tt := []struct {
		name                string
		searchResponse      []models.CloudtrailData
		searchDatabaseError error
		expectedBody        []models.CloudtrailData
		expectedCode        int
	}{
		{
			"success",
			[]models.CloudtrailData{},
			nil,
			[]models.CloudtrailData{},
			http.StatusOK,
		},
		{
			"fail to search db",
			[]models.CloudtrailData{},
			errors.ErrorToRetrieveRecordsFromDb,
			[]models.CloudtrailData{},
			http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		db := database.NewClient()
		db.SetLastRegistersError(tc.searchDatabaseError)
		db.SetLastResgistersResponse(tc.searchResponse)
		a := api{
			db: db,
		}
		r := mux.NewRouter()
		r.Handle("/search", a.Search())

		ts := httptest.NewServer(r)
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/search")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != tc.expectedCode {
			t.Errorf("TEST [%s] FAILED: expected code [%d] but received [%d]", tc.name, tc.expectedCode, resp.StatusCode)
		}
		// rawResp, _ := json.Marshal(resp.Body)
		// expectedRawResp, _ := json.Marshal(tc.expectedBody)
		// if !bytes.Equal(rawResp, expectedRawResp) {
		// t.Errorf("TEST [%s] FAILED: expected body [%s] but received [%s]", tc.name, expectedRawResp, rawResp)
		// }
	}
}

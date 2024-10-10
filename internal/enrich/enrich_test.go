package enrich

import (
	"data-enrich/internal/errors"
	"data-enrich/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCountryFromIp(t *testing.T) {
	tt := []struct {
		name                    string
		externalAPIResponse     any
		externalAPICode         int
		enableExternalServer    bool
		expectedCountryResponse string
		expectedCountryError    error
	}{
		{
			name:                    "success",
			externalAPIResponse:     models.Ip2CountryResponse{CountryName: "Brazil"},
			externalAPICode:         http.StatusOK,
			enableExternalServer:    true,
			expectedCountryResponse: "Brazil",
			expectedCountryError:    nil,
		},
		{
			name:                    "server fail",
			externalAPIResponse:     []byte{},
			externalAPICode:         http.StatusInternalServerError,
			enableExternalServer:    false,
			expectedCountryResponse: "",
			expectedCountryError:    errors.ErrorToRetrieveDataFromUri,
		},
		{
			name:                    "server bad format",
			externalAPIResponse:     `{"bad_json": ""`,
			externalAPICode:         http.StatusOK,
			enableExternalServer:    true,
			expectedCountryResponse: "",
			expectedCountryError:    errors.ErrorFailToReadHttpResponseBody,
		},
	}
	for _, tc := range tt {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.externalAPICode)
			rawResponse, _ := json.Marshal(tc.externalAPIResponse)
			w.Write(rawResponse)
		}))
		if !tc.enableExternalServer {
			s.Close()
		} else {
			defer s.Close()
		}
		enrichService := enrich{
			countryAPIAddress: s.URL + "?"}
		country, err := enrichService.GetCountryFromIp("10.10.10.10")
		if err != tc.expectedCountryError {
			t.Errorf("TEST [%s] FAILED: expected error [%v] but received [%v]", tc.name, tc.expectedCountryError, err)
		}
		if country != tc.expectedCountryResponse {
			t.Errorf("TEST [%s] FAILED: expected country response [%s] but received [%s]", tc.name, tc.expectedCountryResponse, country)
		}
	}
}

func TestGetCountryRegion(t *testing.T) {
	tt := []struct {
		name                   string
		testCountry            string
		externalAPIResponse    any
		externalAPICode        int
		enableExternalServer   bool
		expectedRegionResponse string
		expectedRegionyError   error
	}{
		{
			name:        "success",
			testCountry: "Brazil",
			externalAPIResponse: []models.Country{{Name: models.Name{
				Common: "Brazil",
			},
				Region: "americas",
			}},
			externalAPICode:        http.StatusOK,
			enableExternalServer:   true,
			expectedRegionResponse: "americas",
			expectedRegionyError:   nil,
		},
		{
			name:        "invalid region",
			testCountry: "Antarctica",
			externalAPIResponse: []models.Country{{Name: models.Name{
				Common: "Brazil",
			},
				Region: "South America",
			}},
			externalAPICode:        http.StatusOK,
			enableExternalServer:   true,
			expectedRegionResponse: "",
			expectedRegionyError:   errors.ErrorRegionNotFound,
		},
		{
			name:                   "request error",
			testCountry:            "Brazil",
			externalAPIResponse:    []byte{},
			externalAPICode:        http.StatusInternalServerError,
			enableExternalServer:   false,
			expectedRegionResponse: "",
			expectedRegionyError:   errors.ErrorToRetrieveDataFromUri,
		},
		{
			name:                   "invalid API response",
			testCountry:            "Brazil",
			externalAPIResponse:    `{"bad_json": ""`,
			externalAPICode:        http.StatusInternalServerError,
			enableExternalServer:   true,
			expectedRegionResponse: "",
			expectedRegionyError:   errors.ErrorToUnmarshallRequestBody,
		},
	}
	for _, tc := range tt {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.externalAPICode)
			rawResponse, _ := json.Marshal(tc.externalAPIResponse)
			w.Write(rawResponse)
		}))
		if !tc.enableExternalServer {
			s.Close()
		} else {
			defer s.Close()
		}
		enrichService := enrich{
			regionAPIAddress: s.URL + "/"}
		region, err := enrichService.GetCountryRegion(tc.testCountry)
		if err != tc.expectedRegionyError {
			t.Errorf("TEST [%s] FAILED: expected error [%v] but received [%v]", tc.name, tc.expectedRegionyError, err)
		}
		if region != tc.expectedRegionResponse {
			t.Errorf("TEST [%s] FAILED: expected country response [%s] but received [%s]", tc.name, tc.expectedRegionResponse, region)
		}
	}
}

package enrich

import (
	cte "data-enrich/internal/constants"
	"data-enrich/internal/errors"
	"data-enrich/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Enricher interface {
	GetCountryFromIp(ip string) (string, error)
	GetCountryRegion(countr string) (string, error)
}

type enrich struct {
	countryAPIAddress string
	regionAPIAddress  string
}

func NewEnrichService() Enricher {
	return &enrich{
		countryAPIAddress: cte.IP_API_ADDRESS,
		regionAPIAddress:  cte.REGION_API_ADDRESS,
	}
}

func (e *enrich) GetCountryFromIp(ip string) (string, error) {
	// Make a get request to formatted url
	resp, err := http.Get(e.countryAPIAddress + ip)
	if err != nil {
		log.Print(err)
		return "", errors.ErrorToRetrieveDataFromUri
	}
	defer resp.Body.Close()

	var ip2countryResponse models.Ip2CountryResponse
	err = json.NewDecoder(resp.Body).Decode(&ip2countryResponse)
	if err != nil {
		log.Print(err)
		return "", errors.ErrorFailToReadHttpResponseBody
	}
	return ip2countryResponse.CountryName, nil
}

func (e *enrich) GetCountryRegion(country string) (string, error) {
	// Searches for the country region by region and stops
	// the search if the name of the country is present in the analyzed region
	i := 0
	for {
		regionName := cte.Region(i).String()
		if regionName == "" {
			return "", errors.ErrorRegionNotFound
		}

		// Read http response
		resp, err := http.Get(e.regionAPIAddress + regionName)
		if err != nil {
			log.Print(err)
			return "", errors.ErrorToRetrieveDataFromUri
		}

		defer resp.Body.Close()
		var countryList []models.Country

		// Unmarshal request body
		err = json.NewDecoder(resp.Body).Decode(&countryList)
		if err != nil {
			log.Print(err)
			return "", errors.ErrorToUnmarshallRequestBody
		}

		// Linear search by country name in the region data
		for _, countryInfo := range countryList {
			if country == countryInfo.Name.Common {
				return countryInfo.Region, nil
			}
		}
		i++
	}
}

// Concatenate error messages and return a new error
func WrapError(err error, msg string) error {
	return fmt.Errorf(msg + ": " + err.Error())
}

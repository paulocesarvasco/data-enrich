package utils

import (
	cte "data-enrich/internal/constants"
	"data-enrich/internal/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func GetCountryFromIp(ip string) (string, error) {
	// Make a get request to formatted url
	resp, err := http.Get(cte.IP_API_ADDRESS + ip)
	if err != nil {
		log.Print(err)
		return "", errors.New(cte.ErrorToRetrieveDataFromUri)
	}
	defer resp.Body.Close()

	var ip2countryResponse models.Ip2CountryResponse
	err = json.NewDecoder(resp.Body).Decode(&ip2countryResponse)
	if err != nil {
		log.Print(err)
		return "", errors.New(cte.FailToReadHttpResponseBody)
	}
	return ip2countryResponse.CountryName, nil
}

func GetCountryRegion(country string) (string, error) {
	// Searches for the country region by region and stops
	// the search if the name of the country is present in the analyzed region
	i := 0
	for {
		regionName := cte.Region(i).String()
		if regionName == "" {
			return "", errors.New(cte.RegionNotFound)
		}

		// Read http response
		resp, err := http.Get(cte.REGION_API_ADDRESS + regionName)
		if err != nil {
			log.Print(err)
			return "", errors.New(cte.ErrorToRetrieveDataFromUri)
		}

		defer resp.Body.Close()
		var countryList []models.Country

		// Unmarshal request body
		err = json.NewDecoder(resp.Body).Decode(&countryList)
		if err != nil {
			log.Print(err)
			return "", errors.New(cte.ErrorToUnmarshallRequestBody)
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
	return errors.New(msg + ": " + err.Error())
}

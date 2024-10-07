package utils

import (
	cte "data-enrich/internal/constants"
	"data-enrich/internal/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client
var dbCollection *mongo.Collection

func GetCountryFromIp(ip string) (string, error) {

	url := "https://api.ip2country.info/ip?" + ip

	// Make a get request to formatted url
	resp, err := http.Get(url)
	if err != nil {
		return "", WrapError(err, cte.ErrorToRetrieveDataFromUri)
	}

	defer resp.Body.Close()

	// Read http response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ip2countryResponse models.Ip2CountryResponse

	// Unmarshal request body
	err = json.Unmarshal(body, &ip2countryResponse)
	if err != nil {
		return "", WrapError(err, cte.FailToReadHttpResponseBody)
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
		url := "https://restcountries.com/v3.1/region/" + regionName

		// Read http response
		resp, err := http.Get(url)
		if err != nil {
			return "", WrapError(err, cte.ErrorToRetrieveDataFromUri)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", WrapError(err, cte.FailToReadHttpResponseBody)
		}

		var countryList []models.Country

		// Unmarshal request body
		err = json.Unmarshal(body, &countryList)
		if err != nil {
			return "", WrapError(err, cte.ErrorToUnmarshallRequestBody)
		}

		defer resp.Body.Close()

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

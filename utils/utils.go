package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meli/constants"
	"meli/models"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client
var dbCollection *mongo.Collection

func GetCountryFromIp(ip string) (string, error) {
	url := "https://api.ip2country.info/ip?" + ip
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ip2countryResponse models.Ip2CountryResponse
	// Unmarshal request body
	err = json.Unmarshal(body, &ip2countryResponse)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	return ip2countryResponse.CountryName, nil
}

func GetCountryRegion(country string) (string, error) {

	i := 0
	for {
		regionName := constants.Region(i).String()
		if regionName == "" {
			return "", fmt.Errorf("Regiao nao encontrada")
		}
		url := "https://restcountries.com/v3.1/region/" + regionName

		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var countryList []models.Country
		// Unmarshal request body
		err = json.Unmarshal(body, &countryList)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		for _, countryInfo := range countryList {
			if country == countryInfo.Name.Common {
				return countryInfo.Region, nil
			}
		}

		i++
	}
}

func WrapError(err error, msg string) error {
	return errors.New(msg + ": " + err.Error())
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meli/constants"
	"meli/models"
	"net/http"
	"os/user"
	"path/filepath"
)

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := filepath.Join(dir, "input.json")
	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	var records models.CloudtrailData

	err = json.Unmarshal(file, &records)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug

	}

	sourceIP := records.Records[0].SourceIPAddress

	url := "https://api.ip2country.info/ip?" + sourceIP

func getCountryFromIp(ip string) string {
	url := "https://api.ip2country.info/ip?" + ip
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	var ip2countryResponse models.Ip2CountryResponse
	// Unmarshal request body
	err = json.Unmarshal(body, &ip2countryResponse)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	resp.Body.Close()

	return ip2countryResponse.CountryName
}

func getCountryRegion(country string) string {

	i := 0
	for {
		regionName := constants.Region(i).String()
		if regionName == "" {
			fmt.Printf("%+v\n", "Região não encontrada") // output for debug
			return ""
		}
		url := "https://restcountries.com/v3.1/region/" + regionName

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%+v\n", err) // output for debug
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%+v\n", err) // output for debug
		}

		var countryList []models.Country
		// Unmarshal request body
		err = json.Unmarshal(body, &countryList)
		if err != nil {
			fmt.Printf("%+v\n", err) // output for debug
		}

		defer resp.Body.Close()

		for _, countryInfo := range countryList {
			if country == countryInfo.Name.Common {
				return countryInfo.Region
			}
		}

		i++
	}
}

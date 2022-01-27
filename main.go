package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meli/models"
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
	fmt.Printf("%+v\n", records) // output for debug


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

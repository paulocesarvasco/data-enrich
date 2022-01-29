package connection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"meli/database"
	"meli/models"
	"meli/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/input", saveDataOnDb).Methods("POST").Schemes("http")
	myRouter.HandleFunc("/get", getData).Methods("GET").Schemes("http")
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func saveDataOnDb(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)

	var record models.CloudtrailData

	err := json.Unmarshal(reqBody, &record)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	country, err := utils.GetCountryFromIp(record.Records[0].SourceIPAddress)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	region, err := utils.GetCountryRegion(country)
	if err != nil {
		fmt.Printf("%+v\n", err) // output for debug
	}

	enrichiment := models.Enrichment{
		Country: country,
		Region:  region,
	}

	record.Records[0].Enrichment = enrichiment

	err = database.SaveDataOnDatabase(record)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok\n"))
	}

}

func getData(w http.ResponseWriter, r *http.Request) {

	records, err := database.RetriveLastDataFromDatabase(10)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	byteRecords, err := json.Marshal(records)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(byteRecords)
	}
}

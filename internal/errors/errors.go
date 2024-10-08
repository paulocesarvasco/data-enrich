package errors

import "fmt"

var ErrorDatabaseOperationSave = fmt.Errorf("Failed to store new record in the database")
var ErrorToCreateDatabaseInstace = fmt.Errorf("Failed to create database instance")
var ErrorToCreateDatabaseCollection = fmt.Errorf("Failed to create collection on database")

var ErrorToUnmarshallRequestBody = fmt.Errorf("Failed to decode http request body")
var ErrorToRetriveCountryFromIp = fmt.Errorf("Failed to get country name from external API")
var ErrorToRetriveRegionName = fmt.Errorf("Failed to get country region name")
var ErrortoGetRecordsFromDatabase = fmt.Errorf("Couldn't retrieve data from database")
var ErrortoEncodeDataFromDatabase = fmt.Errorf("Failed to encode data retrieved from database")
var ErrorMissedMandatoryFields = fmt.Errorf("Source IP address necessary to enrich is missing")

var ErrorToEstablishDatabaseConnection = fmt.Errorf("Couldn't establish connection with database service")
var ErrorTestingConnectionWithDB = fmt.Errorf("Database connection not responding")
var ErrorDatabaseClientNotInitialized = fmt.Errorf("Database client not initialized")
var ErrorCollectionNotFound = fmt.Errorf("Couldn't find the collection")
var ErrorToRetrieveRecordsFromDb = fmt.Errorf("Failed when trying to retrieve data from the database")

var ErrorToRetrieveDataFromUri = fmt.Errorf("Unable to retrieve data from URI")
var ErrorFailToReadHttpResponseBody = fmt.Errorf("Couldn't parse data from http response")
var ErrorRegionNotFound = fmt.Errorf("Unknown region")

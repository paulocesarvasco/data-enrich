package database

import (
	"context"
	"errors"
	cte "meli/constants"
	"meli/models"
	"meli/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbClient *mongo.Client
var dbCollection *mongo.Collection
var ctx context.Context

func CreateDatabaseInstance(uri string) error {

	ctx = context.TODO()

	var err error
	dbClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return utils.WrapError(err, cte.ErrorToEstablishDatabaseConnection)
	}

	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return utils.WrapError(err, cte.ErrorTestingConnectionWithDB)
	}

	return nil
}

func CreateDatabaseCollection(dbName string, collectionName string) error {

	if dbClient == nil {
		return errors.New(cte.DatabaseClientNotInitialized)
	}

	dbCollection = dbClient.Database(dbName).Collection(collectionName)

	return nil
}

func SaveDataOnDatabase(data interface{}) error {

	_, err := dbCollection.InsertOne(ctx, data)
	if err != nil {
		return utils.WrapError(err, cte.ErrortoSaveDataOnDatabase)
	}

	return nil
}

func RetriveLastDataFromDatabase(numRegisters int) ([]models.CloudtrailData, error) {

	cur, err := dbCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapError(err, cte.CollectionNotFound)
	}

	defer cur.Close(ctx)

	var allRecords []models.CloudtrailData
	var lastRecords []models.CloudtrailData

	err = cur.All(ctx, &allRecords)
	if err != nil {
		return nil, utils.WrapError(err, cte.ErrorToRetrieveRecordsFromDb)
	}

	if len(allRecords) <= 10 {
		return allRecords, nil
	}

	for i := len(allRecords) - 1; i >= len(allRecords)-numRegisters; i-- {
		lastRecords = append(lastRecords, allRecords[i])
	}

	return lastRecords, nil
}

package database

import (
	"context"
	cte "data-enrich/internal/constants"
	"data-enrich/internal/models"
	"data-enrich/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB interface {
	Save(ctx context.Context, data any) error
	RetriveLastregisters(ctx context.Context, numRegisters int) ([]models.CloudtrailData, error)
}

type db struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewClient() (*db, error) {

	// Initialize db client
	dbClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		return nil, utils.WrapError(err, cte.ErrorToEstablishDatabaseConnection)
	}

	// Test connection between server and client db
	err = dbClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, utils.WrapError(err, cte.ErrorTestingConnectionWithDB)
	}

	dbCollection := dbClient.Database("CloudtrailRecords").Collection("records")

	return &db{
		client:     dbClient,
		collection: dbCollection}, nil

}

// Stores data on previous collection initialized
func (d *db) Save(ctx context.Context, data interface{}) error {
	_, err := d.collection.InsertOne(ctx, data)
	if err != nil {
		return utils.WrapError(err, cte.ErrortoSaveDataOnDatabase)
	}
	return nil
}

// Get last records from database. numRegisters indicates how many records must be retrieved
func (d *db) RetriveLastregisters(ctx context.Context, numRegisters int) ([]models.CloudtrailData, error) {
	cur, err := d.collection.Find(ctx, bson.D{})
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

	// Edge case for when less than 10 records are stored
	if len(allRecords) <= 10 {
		return allRecords, nil
	}

	for i := len(allRecords) - 1; i >= len(allRecords)-numRegisters; i-- {
		lastRecords = append(lastRecords, allRecords[i])
	}

	return lastRecords, nil
}

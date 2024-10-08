package database

import (
	"context"
	cte "data-enrich/internal/constants"
	"data-enrich/internal/models"
	"data-enrich/internal/utils"

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

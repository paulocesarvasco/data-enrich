package database

import (
	"context"
	"data-enrich/internal/errors"
	"data-enrich/internal/models"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB interface {
	Save(ctx context.Context, data any) error
	RetrieveLastRegisters(ctx context.Context, numRegisters int) ([]models.CloudtrailData, error)
}

type db struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewClient() (*db, error) {
	// Initialize db client
	dbClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Print(err)
		return nil, errors.ErrorToEstablishDatabaseConnection
	}

	// Test connection between server and client db
	err = dbClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, errors.ErrorTestingConnectionWithDB
	}
	dbCollection := dbClient.Database("CloudtrailRecords").Collection("records")
	return &db{
		client:     dbClient,
		collection: dbCollection}, nil
}

package database

import (
	"context"
	"fmt"
	"meli/models"

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
		return err
	}

	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

func CreateDatabaseCollection(dbClient *mongo.Client, dbName string, collectionName string) *mongo.Collection {

	dbCollection := dbClient.Database(dbName).Collection(collectionName)

	return dbCollection
}

func PostDataOnDatabase(dbCollection *mongo.Collection, data interface{}) error {

	_, err := dbCollection.InsertOne(context.TODO(), data)
	// check for errors in the insertion
	if err != nil {
		return err
	}

	return nil
}

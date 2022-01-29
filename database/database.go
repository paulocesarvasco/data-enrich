package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateDatabaseInstance(uri string) (*mongo.Client, error) {

	dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = dbClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return dbClient, nil
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

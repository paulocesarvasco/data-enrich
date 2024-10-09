package database

import (
	"context"
	"data-enrich/internal/constants"
	"data-enrich/internal/errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSave(t *testing.T) {
	tt := []struct {
		name           string
		mockResult     primitive.D
		expectedResult error
	}{
		{"success", mtest.CreateSuccessResponse(bson.E{}), nil},
		{"fail", mtest.CreateWriteErrorsResponse(
			mtest.WriteError{Message: constants.ErrorDatabaseOperationSave}), errors.ErrorDatabaseOperationSave},
	}
	for _, tc := range tt {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock).ShareClient(true))
		defer mt.Close()

		mt.AddMockResponses(tc.mockResult)
		dbClient := mt.Client
		dbCollection := dbClient.Database("mongodb://mock").Collection("records")
		db := db{
			client:     dbClient,
			collection: dbCollection}

		data := bson.D{{}}
		err := db.Save(context.Background(), data)
		if err != tc.expectedResult {
			t.Errorf("TEST FAILED: expected result: %s. obtained result: %s", tc.expectedResult, err)
		}
	}
}

// func TestRetriveLastRegisters(t *testing.T) {
//	tt := []struct {
//		name                 string
//		numRegisters         int
//		expectedNumRegisters int
//		expectedError        error
//	}{
//		{"success", 3, 3, nil},
//	}

//	for _, tc := range tt {
//		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock).ShareClient(true).CollectionName("records"))
//		defer mt.Close()

//		dbClient := mt.Client
//		dbCollection := dbClient.Database("mongodb://mock").Collection("records")

//		for i := 0; i < tc.numRegisters; i++ {
//			dbCollection.InsertOne(mtest.Background, bson.E{})
//		}

//		db := db{
//			client:     dbClient,
//			collection: dbCollection}
//		res, err := db.RetriveLastregisters(mtest.Background, tc.numRegisters)
//		if err != tc.expectedError {
//			t.Errorf("TEST FAILED: expected error: %s. obtained result: %s", tc.expectedError, err)
//		}
//		if len(res) != tc.expectedNumRegisters {
//			t.Errorf("TEST FAILED: expected #results: %d. obtained: %d", tc.expectedNumRegisters, len(res))
//		}
//	}
// }

// func TestRetrieveLastRegisters(t *testing.T) {
//	tt := []struct {
//		name                 string
//		numRegisters         int
//		expectedNumRegisters int
//		expectedError        error
//	}{
//		{"success", 3, 3, nil},
//	}

//	for _, tc := range tt {
//		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock).ShareClient(true).CollectionName("records"))
//		defer mt.Close()

//		// Prepare mock responses for InsertOne and Find operations
//		var responses []bson.D
//		for i := 0; i < tc.numRegisters; i++ {
//			// responses = append(responses, bson.D{{Key: "id", Value: i}})
//			responses = append(responses, bson.D{})
//		}

//		mt.AddMockResponses(mtest.CreateSuccessResponse()) // Mock Find close cursor operation
//		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.records", mtest.FirstBatch, responses...))
//		mt.AddMockResponses(mtest.CreateSuccessResponse()) // Mock Find close cursor operation

//		// Setup mocked collection
//		dbClient := mt.Client
//		dbMock := dbClient.Database("mongodb://mock")
//		err := dbMock.CreateCollection(mtest.Background, "records")
//		if err != nil {
//			t.Fatalf("TEST FAILED: failed to configure environment: %v", err)
//		}
//		dbCollection := dbClient.Database("mongodb://mock").Collection("records")

//		// for i := 0; i < tc.numRegisters; i++ {
//		// _, err := dbCollection.InsertOne(mtest.Background, bson.E{})
//		// if err != nil {
//		// t.Fatalf("Failed to insert document: %v", err)
//		// }
//		// }

//		db := db{
//			client:     dbClient,
//			collection: dbCollection,
//		}
//		res, err := db.RetriveLastregisters(mtest.Background, tc.numRegisters)
//		if err != tc.expectedError {
//			t.Errorf("TEST FAILED: expected error: %v, got: %v", tc.expectedError, err)
//		}
//		if len(res) != tc.expectedNumRegisters {
//			t.Errorf("TEST FAILED: expected #results: %d, got: %d", tc.expectedNumRegisters, len(res))
//		}
//	}
// }

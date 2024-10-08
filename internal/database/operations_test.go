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

		data := bson.D{{Key: "username", Value: "duplicate_user"}}

		err := db.Save(context.Background(), data)
		if err != tc.expectedResult {
			t.Errorf("TEST FAILED: expected result: %s. obtained result: %s", tc.expectedResult, err)
		}
	}
}

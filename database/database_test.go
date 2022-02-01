package database

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbInstanceTest struct {
	uri           string
	errorExpected string
}

var instanceTests = []dbInstanceTest{
	{"", "error parsing uri: scheme must be \"mongodb\" or \"mongodb+srv\""},
	{"mongodb", "error parsing uri: scheme must be \"mongodb\" or \"mongodb+srv\""},
}

func TestCreateDatabaseInstance(t *testing.T) {

	ctx := context.TODO()

	for _, test := range instanceTests {
		_, err := mongo.Connect(ctx, options.Client().ApplyURI(test.uri))
		if err.Error() != test.errorExpected {
			t.Errorf("%+v != %+v", err.Error(), test.errorExpected)
		}
	}
}

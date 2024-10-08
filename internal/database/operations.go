package database

import (
	"context"
	cte "data-enrich/internal/constants"
	"data-enrich/internal/models"
	"data-enrich/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// Stores data on previous collection initialized
func (d *db) Save(ctx context.Context, data any) error {
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

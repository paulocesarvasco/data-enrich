package mocks

import (
	"context"
	"data-enrich/internal/models"
)

type dbClient struct {
	saveError             error
	lastRegistersError    error
	lastRegistersResponse []models.CloudtrailData
}

func NewClient() *dbClient {
	return &dbClient{}
}

func (db *dbClient) SetSaveError(err error) {
	db.saveError = err
}

func (db *dbClient) SetLastResgistersResponse(res []models.CloudtrailData) {
	db.lastRegistersResponse = res
}

func (db *dbClient) SetLastRegistersError(err error) {
	db.lastRegistersError = err
}

func (db *dbClient) Save(ctx context.Context, data any) error {
	return db.saveError
}

func (db *dbClient) RetrieveLastRegisters(ctx context.Context, numRegisters int) ([]models.CloudtrailData, error) {
	return db.lastRegistersResponse, db.lastRegistersError
}

package repository

import (
	"bank-statement/models"
	"context"
)

type StatementRepositoryInterface interface {
	Find(ctx context.Context, filter map[string]interface{}) (result []models.Statement, err error)
	BulkCreate(ctx context.Context, statements []models.Statement) (err error)
	Summary(ctx context.Context) (result models.StatementSummary, err error)
}

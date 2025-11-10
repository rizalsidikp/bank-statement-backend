package repository

import (
	"bank-statement/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type StatementRepositoryMock struct {
	Mock mock.Mock
}

func (m *StatementRepositoryMock) Find(ctx context.Context, filter map[string]interface{}) (result []models.Statement, err error) {
	args := m.Mock.Called(ctx, filter)
	return args.Get(0).([]models.Statement), args.Error(1)
}

func (m *StatementRepositoryMock) BulkCreate(ctx context.Context, statements []models.Statement) (err error) {
	args := m.Mock.Called(ctx, statements)
	return args.Error(0)
}

func (m *StatementRepositoryMock) Summary(ctx context.Context) (result models.StatementSummary, err error) {
	args := m.Mock.Called(ctx)
	return args.Get(0).(models.StatementSummary), args.Error(1)
}

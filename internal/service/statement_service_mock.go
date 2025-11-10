package service

import (
	"bank-statement/dto"
	"context"

	"github.com/stretchr/testify/mock"
)

type StatementServiceMock struct {
	Mock mock.Mock
}

func (m *StatementServiceMock) ListIssuedStatements(ctx context.Context) (result []dto.StatementDTO, err error) {
	args := m.Mock.Called(ctx)
	return args.Get(0).([]dto.StatementDTO), args.Error(1)
}

func (m *StatementServiceMock) BulkCreateStatements(ctx context.Context, dtos []dto.StatementDTO) (err error) {
	args := m.Mock.Called(ctx, dtos)
	return args.Error(0)
}

func (m *StatementServiceMock) CalculateStatementBalance(ctx context.Context) (result dto.StatementSummaryDTO, err error) {
	args := m.Mock.Called(ctx)
	return args.Get(0).(dto.StatementSummaryDTO), args.Error(1)
}

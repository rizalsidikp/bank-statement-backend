package service

import (
	"bank-statement/dto"
	"context"
)

type StatementServiceInterface interface {
	ListIssuedStatements(ctx context.Context) (result []dto.StatementDTO, err error)
	BulkCreateStatements(ctx context.Context, dtos []dto.StatementDTO) (err error)
	CalculateStatementBalance(ctx context.Context) (result dto.StatementSummaryDTO, err error)
}

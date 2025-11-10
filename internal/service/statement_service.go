package service

import (
	"bank-statement/dto"
	"bank-statement/internal/repository"
	"context"
)

type StatementService struct {
	repository repository.StatementRepositoryInterface
}

func NewStatementService(repository repository.StatementRepositoryInterface) StatementServiceInterface {
	return &StatementService{
		repository: repository,
	}
}

func (s *StatementService) ListIssuedStatements(ctx context.Context) (result []dto.StatementDTO, err error) {
	statements, err := s.repository.Find(ctx, map[string]interface{}{
		"status": []string{"FAILED", "PENDING"},
	})
	if err != nil {
		return nil, err
	}
	result = dto.FromStatementModels(statements)
	return result, nil
}

func (s *StatementService) BulkCreateStatements(ctx context.Context, dtos []dto.StatementDTO) (err error) {
	statements := dto.ToStatementModels(dtos)
	return s.repository.BulkCreate(ctx, statements)
}

func (s *StatementService) CalculateStatementBalance(ctx context.Context) (result dto.StatementSummaryDTO, err error) {
	summary, err := s.repository.Summary(ctx)
	if err != nil {
		return result, err
	}
	result.FromSummaryModel(&summary)
	return result, nil
}

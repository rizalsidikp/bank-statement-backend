package service

import (
	"bank-statement/dto"
	"bank-statement/internal/repository"
	"bank-statement/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var statementRepository = repository.StatementRepositoryMock{Mock: mock.Mock{}}
var statementService = NewStatementService(&statementRepository)

func TestListIssuedStatements(t *testing.T) {
	t.Run("should return issued statements", func(t *testing.T) {
		statementRepository.Mock.On("Find", mock.Anything, map[string]interface{}{
			"status": []string{"FAILED", "PENDING"},
		}).Return([]models.Statement{
			{
				ID:          "4768b7e4-de41-4e26-a61f-ef127615b168",
				Timestamp:   1624608050,
				Name:        "E-COMMERCE A",
				Type:        "DEBIT",
				Amount:      150000,
				Status:      "FAILED",
				Description: "clothes",
			},
		}, nil).Once()

		result, err := statementService.ListIssuedStatements(context.Background())
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "E-COMMERCE A", result[0].Name)
		assert.Equal(t, int64(150000), result[0].Amount)
		assert.Equal(t, models.FAILED, result[0].Status)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		statementRepository.Mock.On("Find", mock.Anything, map[string]interface{}{
			"status": []string{"FAILED", "PENDING"},
		}).Return([]models.Statement{}, assert.AnError).Once()
		result, err := statementService.ListIssuedStatements(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBulkCreateStatements(t *testing.T) {
	t.Run("should create statements successfully", func(t *testing.T) {
		dtos := []dto.StatementDTO{
			{
				Timestamp:   1624608080,
				Name:        "RESTAURANT B",
				Type:        "DEBIT",
				Amount:      200000,
				Status:      "PENDING",
				Description: "dinner",
			},
		}
		statements := dto.ToStatementModels(dtos)
		statementRepository.Mock.On("BulkCreate", mock.Anything, statements).Return(nil).Once()
		err := statementService.BulkCreateStatements(context.Background(), dtos)
		assert.NoError(t, err)
	})
}

func TestCalculateStatementBalance(t *testing.T) {
	t.Run("should return correct statement balance", func(t *testing.T) {
		statementRepository.Mock.On("Summary", mock.Anything).Return(models.StatementSummary{
			TotalDebit:  300000,
			TotalCredit: 5000000,
		}, nil).Once()
		result, err := statementService.CalculateStatementBalance(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(300000), result.TotalDebit)
		assert.Equal(t, int64(5000000), result.TotalCredit)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		statementRepository.Mock.On("Summary", mock.Anything).Return(models.StatementSummary{}, assert.AnError).Once()
		result, err := statementService.CalculateStatementBalance(context.Background())
		assert.Error(t, err)
		assert.Equal(t, dto.StatementSummaryDTO{}, result)
	})
}

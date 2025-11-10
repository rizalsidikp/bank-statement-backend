package repository

import (
	"bank-statement/database"
	"bank-statement/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db = database.Database{
	Statements: []models.Statement{
		{
			ID:          "4768b7e4-de41-4e26-a61f-ef127615b168",
			Timestamp:   1624608050,
			Name:        "E-COMMERCE A",
			Type:        "DEBIT",
			Amount:      150000,
			Status:      "FAILED",
			Description: "clothes",
		},
		{
			ID:          "5768b7e4-de41-4e26-a61f-ef127615b169",
			Timestamp:   1624608060,
			Name:        "SALARY JUNE",
			Type:        "CREDIT",
			Amount:      5000000,
			Status:      "SUCCESS",
			Description: "monthly salary",
		},
		{
			ID:          "6768b7e4-de41-4e26-a61f-ef127615b170",
			Timestamp:   1624608070,
			Name:        "GROCERY STORE",
			Type:        "DEBIT",
			Amount:      300000,
			Status:      "SUCCESS",
			Description: "weekly groceries",
		},
	},
}

func TestFind(t *testing.T) {
	repo := NewStatementRepository(db)

	t.Run("should return statements matching the filter [status: FAILED]", func(t *testing.T) {
		filter := map[string]interface{}{
			"status": "FAILED",
		}
		statements, err := repo.Find(context.Background(), filter)
		assert.NoError(t, err)
		assert.Len(t, statements, 1)
		assert.Equal(t, "E-COMMERCE A", statements[0].Name)
		assert.Equal(t, int64(150000), statements[0].Amount)
		assert.Equal(t, models.FAILED, statements[0].Status)
	})

	t.Run("should return statements matching the filter [status: FAILED, PENDING]", func(t *testing.T) {
		filter := map[string]interface{}{
			"status": []string{"FAILED", "PENDING"},
		}
		statements, err := repo.Find(context.Background(), filter)
		assert.NoError(t, err)
		assert.Len(t, statements, 1)
		assert.Equal(t, "E-COMMERCE A", statements[0].Name)
		assert.Equal(t, int64(150000), statements[0].Amount)
		assert.Equal(t, models.FAILED, statements[0].Status)
	})

	t.Run("should return statements matching the filter integer", func(t *testing.T) {
		filter := map[string]interface{}{
			"status": 123,
		}
		statements, err := repo.Find(context.Background(), filter)
		assert.NoError(t, err)
		assert.Len(t, statements, 0)
	})
}

func TestSummary(t *testing.T) {
	repo := NewStatementRepository(db)
	t.Run("should return correct summary", func(t *testing.T) {
		summary, err := repo.Summary(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(300000), summary.TotalDebit)
		assert.Equal(t, int64(5000000), summary.TotalCredit)
	})
}

func TestInsert(t *testing.T) {
	repo := NewStatementRepository(db)
	t.Run("should insert statements successfully", func(t *testing.T) {
		statements := []models.Statement{
			{
				Timestamp:   1624608070,
				Name:        "E-COMMERCE B",
				Type:        "DEBIT",
				Amount:      250000,
				Status:      "SUCCESS",
				Description: "electronics",
			},
			{
				Timestamp:   1624608080,
				Name:        "FREELANCE PROJECT",
				Type:        "CREDIT",
				Amount:      2000000,
				Status:      "PENDING",
				Description: "website development",
			},
		}
		err := repo.BulkCreate(context.Background(), statements)
		assert.NoError(t, err)
		// Verify that the statements were inserted
		allStatements, err := repo.Find(context.Background(), map[string]interface{}{})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(allStatements), 4) // original 2 + 2 new
	})
}

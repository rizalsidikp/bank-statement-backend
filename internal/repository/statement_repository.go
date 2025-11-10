package repository

import (
	"bank-statement/database"
	"bank-statement/models"
	"context"

	"github.com/google/uuid"
)

type StatementRepository struct {
	// db is from package database
	db database.Database
}

func NewStatementRepository(db database.Database) StatementRepositoryInterface {
	return &StatementRepository{
		db: db,
	}
}

func (r *StatementRepository) Find(ctx context.Context, filter map[string]interface{}) (result []models.Statement, err error) {
	for _, statement := range r.db.Statements {
		if matchFilter(statement, filter) {
			result = append(result, statement)
		}
	}
	return result, nil
}

func (r *StatementRepository) BulkCreate(ctx context.Context, statements []models.Statement) (err error) {
	for _, statement := range statements {
		// generate ID
		statement.ID = uuid.New().String()
		r.db.Statements = append(r.db.Statements, statement)
	}
	return nil
}

func (r *StatementRepository) Summary(ctx context.Context) (result models.StatementSummary, err error) {
	var totalDebit int64 = 0
	var totalCredit int64 = 0
	for _, statement := range r.db.Statements {
		if statement.Status != models.SUCCESS {
			continue
		}
		switch statement.Type {
		case models.DEBIT:
			totalDebit += statement.Amount
		case models.CREDIT:
			totalCredit += statement.Amount
		}
	}
	result = models.StatementSummary{
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
	}
	return result, nil
}

func matchFilter(statement models.Statement, filter map[string]interface{}) bool {
	// example filter is Returns non-successful transactions (FAILED + PENDING)
	for key, value := range filter {
		switch key {
		case "status":
			switch v := value.(type) {
			case string:
				if statement.Status != models.StatementStatus(v) {
					return false
				}
			case []string:
				match := false
				for _, s := range v {
					if string(statement.Status) == s {
						match = true
						break
					}
				}
				if !match {
					return false
				}
			default:
				return false
			}
		}
	}
	return true
}

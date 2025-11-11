package dto

import "bank-statement/models"

type StatementDTO struct {
	ID          string                 `json:"id"`
	Timestamp   int64                  `json:"timestamp"`
	Name        string                 `json:"name"`
	Type        models.StatementType   `json:"type"`
	Amount      int64                  `json:"amount"`
	Status      models.StatementStatus `json:"status"`
	Description string                 `json:"description"`
}

type StatementSummaryDTO struct {
	TotalDebit   int64 `json:"total_debit"`
	TotalCredit  int64 `json:"total_credit"`
	TotalBalance int64 `json:"total_balance"`
}

type CreateStatementDTO struct {
	File string `json:"file" binding:"required"`
}

func (dto *StatementDTO) ToStatementModel(statement *models.Statement) {
	statement.ID = dto.ID
	statement.Timestamp = dto.Timestamp
	statement.Name = dto.Name
	statement.Type = dto.Type
	statement.Amount = dto.Amount
	statement.Status = dto.Status
	statement.Description = dto.Description
}

func FromStatementModels(statement []models.Statement) (dtos []StatementDTO) {
	for _, s := range statement {
		dtos = append(dtos, StatementDTO{
			ID:          s.ID,
			Timestamp:   s.Timestamp,
			Name:        s.Name,
			Type:        s.Type,
			Amount:      s.Amount,
			Status:      s.Status,
			Description: s.Description,
		})
	}
	return dtos
}

func ToStatementModels(dtos []StatementDTO) (statements []models.Statement) {
	for _, dto := range dtos {
		statements = append(statements, models.Statement{
			ID:          dto.ID,
			Timestamp:   dto.Timestamp,
			Name:        dto.Name,
			Type:        dto.Type,
			Amount:      dto.Amount,
			Status:      dto.Status,
			Description: dto.Description,
		})
	}
	return statements
}

func (dto *StatementSummaryDTO) FromSummaryModel(summary *models.StatementSummary) {
	dto.TotalDebit = summary.TotalDebit
	dto.TotalCredit = summary.TotalCredit
	dto.TotalBalance = summary.TotalCredit - summary.TotalDebit
}

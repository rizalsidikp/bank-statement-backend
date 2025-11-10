package models

type StatementType string
type StatementStatus string

const (
	DEBIT  StatementType = "DEBIT"
	CREDIT StatementType = "CREDIT"
)

const (
	SUCCESS StatementStatus = "SUCCESS"
	FAILED  StatementStatus = "FAILED"
	PENDING StatementStatus = "PENDING"
)

type Statement struct {
	ID          string          `json:"id"`
	Timestamp   int64           `json:"timestamp"`
	Name        string          `json:"name"`
	Type        StatementType   `json:"type"`
	Amount      int64           `json:"amount"`
	Status      StatementStatus `json:"status"`
	Description string          `json:"description"`
}

type StatementSummary struct {
	TotalDebit  int64 `json:"total_debit"`
	TotalCredit int64 `json:"total_credit"`
}

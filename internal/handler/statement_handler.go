package handler

import (
	"bank-statement/dto"
	"bank-statement/internal/service"
	"bank-statement/models"
	"encoding/csv"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type StatementHandler struct {
	statementService service.StatementServiceInterface
}

func NewStatementHandler(statementService service.StatementServiceInterface) *StatementHandler {
	return &StatementHandler{
		statementService: statementService,
	}
}

func (h *StatementHandler) Routes(route fiber.Router) {
	route.Post("/upload", h.UploadStatement)
	route.Get("/balance", h.GetStatementBalance)
	route.Get("/issues", h.GetStatementIssues)
}

func (h *StatementHandler) UploadStatement(c *fiber.Ctx) error {
	var (
		statusCode int = fiber.StatusOK
		response   dto.Response
		ctx        = c.Context()
	)

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		statusCode = fiber.StatusBadRequest
		response.Message = "Failed to get uploaded file: " + err.Error()
		response.Code = statusCode
		return c.Status(statusCode).JSON(response)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".csv" {
		statusCode = fiber.StatusBadRequest
		response.Message = "Invalid file type. Only CSV files are allowed."
		response.Code = statusCode
		return c.Status(statusCode).JSON(response)
	}

	f, err := file.Open()
	if err != nil {
		statusCode = fiber.StatusInternalServerError
		response.Message = "Failed to open uploaded file: " + err.Error()
		response.Code = statusCode
		return c.Status(statusCode).JSON(response)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	dtos := []dto.StatementDTO{}
	count := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) < 6 {
			continue
		}

		amount, err := strconv.ParseInt(strings.TrimSpace(record[3]), 10, 64)
		if err != nil {
			log.Printf("Skipping record due to invalid amount: %v", record)
			continue
		}

		timestamp, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			log.Printf("Skipping record due to invalid timestamp: %v", record)
			continue
		}

		statement := dto.StatementDTO{
			Timestamp:   int64(timestamp),
			Name:        strings.TrimSpace(record[1]),
			Type:        models.StatementType(strings.TrimSpace(record[2])),
			Amount:      int64(amount),
			Status:      models.StatementStatus(strings.TrimSpace(record[4])),
			Description: strings.TrimSpace(record[5]),
		}
		dtos = append(dtos, statement)
		count++
	}
	err = h.statementService.BulkCreateStatements(ctx, dtos)
	if err != nil {
		statusCode = fiber.StatusInternalServerError
		response.Message = "Failed to create statements: " + err.Error()
		response.Code = statusCode
		return c.Status(statusCode).JSON(response)
	}

	response.Message = "Successfully uploaded " + strconv.Itoa(count) + " statements"
	response.Code = statusCode
	return c.Status(statusCode).JSON(response)
}

func (h *StatementHandler) GetStatementBalance(c *fiber.Ctx) error {
	var (
		statusCode int = fiber.StatusOK
		response   dto.Response
		ctx        = c.Context()
	)

	balance, err := h.statementService.CalculateStatementBalance(ctx)
	if err != nil {
		statusCode = fiber.StatusInternalServerError
		response.Message = "Failed to get statement balance: " + err.Error()
		response.Code = statusCode
		return c.Status(statusCode).JSON(response)
	}

	response.Data = balance
	response.Message = "Successfully retrieved statement balance"
	response.Code = statusCode
	return c.Status(statusCode).JSON(response)
}

func (h *StatementHandler) GetStatementIssues(c *fiber.Ctx) error {
	var (
		statusCode int = fiber.StatusOK
		response   dto.Response
		ctx        = c.Context()
	)

	issues, err := h.statementService.ListIssuedStatements(ctx)
	if err != nil {
		statusCode = fiber.StatusInternalServerError
		response.Message = "Failed to list issued statements: " + err.Error()
		response.Code = statusCode
		return c.Status(statusCode).JSON(response)
	}

	response.Data = issues
	response.Message = "Successfully retrieved statement issues"
	response.Code = statusCode
	return c.Status(statusCode).JSON(response)
}

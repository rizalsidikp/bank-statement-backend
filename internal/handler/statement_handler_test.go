package handler

import (
	"bank-statement/dto"
	"bank-statement/internal/service"
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var statementService = service.StatementServiceMock{Mock: mock.Mock{}}
var statementHandler = NewStatementHandler(&statementService)

func setupTestApp() (*fiber.App, fiber.Router) {
	app := fiber.New()
	router := app.Group("/api/v1")

	return app, router
}

func TestRoutes(t *testing.T) {
	app, _ := setupTestApp()
	statementHandler.Routes(app.Group(""))
}

func TestUploadStatement(t *testing.T) {
	app, router := setupTestApp()

	t.Run("should upload statement successfully", func(t *testing.T) {
		statementService.Mock.On("BulkCreateStatements", mock.Anything, mock.Anything).Return(nil).Once()
		statementHandler.Routes(router)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.csv")
		part.Write([]byte("1624507883, JOHN DOE, DEBIT, 250000, SUCCESS, restaurant\n1624608050, E-COMMERCE A, DEBIT, 150000, FAILED, clothes\n1624512883, COMPANY A, CREDIT, 12000000, SUCCESS, salary\n1624615065, E-COMMERCE B, DEBIT, 150000, PENDING, clothes"))
		writer.Close()

		req := httptest.NewRequest("POST", "/api/v1/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error for invalid file type", func(t *testing.T) {
		statementHandler.Routes(router)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("This is a test text file."))
		writer.Close()
		req := httptest.NewRequest("POST", "/api/v1/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error when file is missing", func(t *testing.T) {
		statementHandler.Routes(router)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()
		req := httptest.NewRequest("POST", "/api/v1/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should upload statement successfully", func(t *testing.T) {
		statementService.Mock.On("BulkCreateStatements", mock.Anything, mock.Anything).Return(nil).Once()
		statementHandler.Routes(router)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.csv")
		part.Write([]byte("1624507883, JOHN DOE, DEBIT, 250000, SUCCESS\n162460805a, E-COMMERCE A, DEBIT, 150000, FAILED, clothes\n1624512883, COMPANY A, CREDIT, 12a00000, SUCCESS, salary\n1624615065, E-COMMERCE B, DEBIT, 150000, PENDING, clothes"))
		writer.Close()

		req := httptest.NewRequest("POST", "/api/v1/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should failed to upload statement due to service error", func(t *testing.T) {
		statementService.Mock.On("BulkCreateStatements", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		statementHandler.Routes(router)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.csv")
		part.Write([]byte("1624507883, JOHN DOE, DEBIT, 250000, SUCCESS, restaurant\n1624608050, E-COMMERCE A, DEBIT, 150000, FAILED, clothes\n1624512883, COMPANY A, CREDIT, 12000000, SUCCESS, salary\n1624615065, E-COMMERCE B, DEBIT, 150000, PENDING, clothes"))
		writer.Close()

		req := httptest.NewRequest("POST", "/api/v1/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestGetStatementBalance(t *testing.T) {
	app, router := setupTestApp()
	t.Run("should return statement balance successfully", func(t *testing.T) {
		statementService.Mock.On("CalculateStatementBalance", mock.Anything).Return(dto.StatementSummaryDTO{
			TotalDebit:  300000,
			TotalCredit: 5000000,
		}, nil).Once()
		statementHandler.Routes(router)
		req := httptest.NewRequest("GET", "/api/v1/balance", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when service fails", func(t *testing.T) {
		statementService.Mock.On("CalculateStatementBalance", mock.Anything).Return(dto.StatementSummaryDTO{}, assert.AnError).Once()
		statementHandler.Routes(router)
		req := httptest.NewRequest("GET", "/api/v1/balance", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestGetStatementIssues(t *testing.T) {
	app, router := setupTestApp()
	t.Run("should return statement issues successfully", func(t *testing.T) {
		statementService.Mock.On("ListIssuedStatements", mock.Anything).Return([]dto.StatementDTO{
			{
				Timestamp:   1624608050,
				Name:        "E-COMMERCE A",
				Type:        "DEBIT",
				Amount:      150000,
				Status:      "FAILED",
				Description: "clothes",
			},
		}, nil).Once()
		statementHandler.Routes(router)
		req := httptest.NewRequest("GET", "/api/v1/issues", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when service fails", func(t *testing.T) {
		statementService.Mock.On("ListIssuedStatements", mock.Anything).Return([]dto.StatementDTO{}, assert.AnError).Once()
		statementHandler.Routes(router)
		req := httptest.NewRequest("GET", "/api/v1/issues", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

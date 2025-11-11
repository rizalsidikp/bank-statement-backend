package main

import (
	"bank-statement/database"
	"bank-statement/internal/handler"
	"bank-statement/internal/repository"
	"bank-statement/internal/service"
	"bank-statement/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()
	// database
	var db database.Database
	db.Statements = []models.Statement{}

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	routeGroup := app.Group("/api/v1")
	routeGroup.Use(recover.New())

	// repository
	statementRepository := repository.NewStatementRepository(db)

	// service
	statementService := service.NewStatementService(statementRepository)

	// handler
	statementHandler := handler.NewStatementHandler(statementService)

	// routes
	statementHandler.Routes(routeGroup)

	// listen and serve on
	app.Listen(":8080")
}

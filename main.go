package main

import (
	"code-review/database"
	"code-review/logger"
	"code-review/models"
	"code-review/routes"
	"github.com/gofiber/fiber/v2"
	"log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	requestLogger "github.com/gofiber/fiber/v2/middleware/logger"
	fRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(requestLogger.New())
	app.Use(fRecover.New())

	// initialize logger
	logger.InitLogger()

	// set up database configurations
	databaseConfig := new(database.DBConnectConfig)
	databaseConfig.MigrationModels = []any{
		&models.User{},
		&models.GHToken{},
		&models.LoginInitSession{},
	}

	databaseConfig.MakeMigrations = true

	// connect to database
	database.ConnectToDB(databaseConfig)

	// setup routes for services
	routes.SetupRoutes(app, []routes.RouteConfig{
		routes.GHRouteConfig,
		routes.AuthRouteConfig,
	})

	log.Panic(app.Listen(":8080"))
}

package main

import (
	"code-review/database"
	"code-review/models"
	"code-review/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(fRecover.New())

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

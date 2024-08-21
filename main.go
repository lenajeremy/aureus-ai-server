package main

import (
	"code-review/auth"
	"code-review/database"
	"code-review/github"
	"code-review/globals"
	"code-review/structs"
	"code-review/utils"
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

	database.ConnectToDB()

	if err := database.DB.AutoMigrate(auth.User{}, globals.BaseModel{}, github.Token{}, globals.LoginInitSession{}); err != nil {
		log.Fatalf(err.Error())
	}

	utils.SetupRoutes(app, []structs.RouteConfig{
		auth.RouteConfig,
		github.RouteConfig,
	})

	log.Panic(app.Listen(":8080"))
}

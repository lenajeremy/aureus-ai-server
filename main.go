package main

import (
	"code-review/auth"
	"code-review/database"
	"code-review/github"
	"code-review/structs"
	"code-review/utils"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()

	database.ConnectToDB()

	utils.SetupRoutes(app, []structs.RouteConfig{
		auth.RouteConfig,
		github.RouteConfig,
	})

	log.Fatal(app.Listen(":8080"))
}

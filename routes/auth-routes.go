package routes

import (
	"code-review/controllers"
	"github.com/gofiber/fiber/v2"
)

var routes = []AppRoute{
	{"/github/initiate", fiber.MethodGet, controllers.InitiateGHLogin},
	{"/github/callback", fiber.MethodGet, controllers.GHLoginCallback},
}

var AuthRouteConfig = RouteConfig{
	Routes:       routes,
	RequiresAuth: false,
	Middleware:   []fiber.Handler{},
	BaseURL:      "auth",
}

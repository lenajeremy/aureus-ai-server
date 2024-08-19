package auth

import (
	"code-review/structs"
	"github.com/gofiber/fiber/v2"
)

var routes = []structs.AppRoute{
	{"/github/initiate", fiber.MethodGet, InitiateGHLogin},
	{"/github/callback", fiber.MethodGet, GHLoginCallback},
}

var RouteConfig = structs.RouteConfig{
	Routes:       routes,
	RequiresAuth: false,
	Middleware:   []fiber.Handler{},
	BaseURL:      "auth",
}

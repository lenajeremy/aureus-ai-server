package auth

import (
	"code-review/structs"
	"github.com/gofiber/fiber/v2"
)

var routes = []structs.AppRoute{
	{"/login", fiber.MethodGet, Login},
	{"/register", fiber.MethodGet, Register},
}

var RouteConfig = structs.RouteConfig{
	Routes:       routes,
	RequiresAuth: false,
	Middleware:   []fiber.Handler{},
	BaseURL:      "auth",
}

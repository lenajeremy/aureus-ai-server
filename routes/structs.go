package routes

import "github.com/gofiber/fiber/v2"

type AppRoute struct {
	Path    string
	Method  string
	Handler fiber.Handler
}

type RouteConfig struct {
	BaseURL      string
	Routes       []AppRoute
	Middleware   []fiber.Handler
	RequiresAuth bool
}

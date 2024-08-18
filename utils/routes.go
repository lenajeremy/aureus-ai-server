package utils

import (
	"code-review/structs"
	"github.com/gofiber/fiber/v2"
)

type RouteFunctionMethod = func(path string, handlers ...fiber.Handler) fiber.Router

func SetupRoute(app *fiber.App, config structs.RouteConfig) {
	router := app.Group(config.BaseURL)

	for _, route := range config.Routes {
		requestMethod := getRequestMethod(router, route.Method)
		requestMethod(route.Path, route.Handler)
	}
}

func SetupRoutes(app *fiber.App, routeConfigs []structs.RouteConfig) {
	for _, config := range routeConfigs {
		SetupRoute(app, config)
	}
}

func getRequestMethod(router fiber.Router, method string) RouteFunctionMethod {
	switch method {
	case fiber.MethodGet:
		return router.Get
	case fiber.MethodPost:
		return router.Post
	case fiber.MethodPut:
		return router.Put
	case fiber.MethodDelete:
		return router.Delete
	case fiber.MethodHead:
		return router.Head
	case fiber.MethodOptions:
		return router.Options
	case fiber.MethodPatch:
		return router.Patch
	case fiber.MethodConnect:
		return router.Connect
	case fiber.MethodTrace:
		return router.Trace
	default:
		return router.Get
	}
}

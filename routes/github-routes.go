package routes

import (
	"code-review/controllers"
	"github.com/gofiber/fiber/v2"
)

var GHRouteConfig = RouteConfig{
	BaseURL:      "github",
	RequiresAuth: false,
	Middleware:   make([]fiber.Handler, 0),

	Routes: []AppRoute{
		{Path: "/pr", Method: fiber.MethodGet, Handler: controllers.HandlePullRequests},
		{Path: "/issues", Method: fiber.MethodGet, Handler: controllers.HandleIssues},
	},
}

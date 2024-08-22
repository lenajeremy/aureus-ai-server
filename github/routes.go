package github

import (
	"code-review/structs"
	"github.com/gofiber/fiber/v2"
)

var RouteConfig = structs.RouteConfig{
	BaseURL:      "/github",
	RequiresAuth: false,
	Middleware:   make([]fiber.Handler, 0),

	Routes: []structs.AppRoute{
		{Path: "/pr", Method: fiber.MethodGet, Handler: HandlePullRequests},
		{Path: "/issues", Method: fiber.MethodGet, Handler: HandleIssues},
	},
}

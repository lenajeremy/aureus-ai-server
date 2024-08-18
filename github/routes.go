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
		{"/pr", fiber.MethodGet, HandlePullRequests},
		{"/issues", fiber.MethodGet, HandleIssues},
	},
}

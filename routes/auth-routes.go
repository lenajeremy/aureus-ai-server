package routes

import (
	"code-review/controllers"
	"github.com/gofiber/fiber/v2"
)

var routes = []AppRoute{
	{"/github/initiate", fiber.MethodGet, controllers.InitiateGHLogin},
	{"/github/callback", fiber.MethodGet, controllers.GHLoginCallback},
	{"/email/signin", fiber.MethodPost, controllers.EmailSignIn},
	{"/email/signup", fiber.MethodPost, controllers.EmailSignUp},
	{"/email/verify", fiber.MethodGet, controllers.VerifyEmail},
}

var AuthRouteConfig = RouteConfig{
	Routes:       routes,
	RequiresAuth: false,
	Middleware:   []fiber.Handler{},
	BaseURL:      "auth",
}

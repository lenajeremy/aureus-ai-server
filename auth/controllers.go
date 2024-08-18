package auth

import "github.com/gofiber/fiber/v2"

func Login(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"hello": "Loginnnn babyyyyy",
	})
}

func Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"hello": "registerrr babyyyyyyy",
	})
}

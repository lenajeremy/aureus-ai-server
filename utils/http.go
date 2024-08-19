package utils

import "github.com/gofiber/fiber/v2"

func RespondSuccess(c *fiber.Ctx, data interface{}, msg string) error {
	return c.JSON(fiber.Map{
		"data":    data,
		"success": true,
		"message": msg,
	})
}

func RespondError(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"data":    nil,
		"success": false,
		"message": err.Error(),
	})
}

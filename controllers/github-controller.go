package controllers

import "github.com/gofiber/fiber/v2"

func HandlePullRequests(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"hello": "handling pull request",
	})
}

func HandleIssues(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"hello": "Handling issuessssss....",
	})
}

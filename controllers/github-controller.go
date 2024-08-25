package controllers

import (
	"code-review/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

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

func GetUserRepos(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)

	if err != nil {
		return utils.RespondError(c, fiber.StatusUnauthorized, err)
	}

	log.Println(user)

	repos, err := utils.GetUserRepos(user.GithubToken.AccessToken)

	if err != nil {
		log.Println(err)
		return utils.RespondError(c, fiber.StatusInternalServerError, errors.New("unable to get user's repositories"))
	}

	return utils.RespondSuccess(c, repos, "Successfully retried user repos")
}

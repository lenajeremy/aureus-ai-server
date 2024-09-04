package controllers

import (
	"code-review/database"
	"code-review/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func HandlePullRequests(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"hello": "handling pull request",
	})
}

func HandleIssues(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"hello": "Handling issue....",
	})
}

func GetUserRepos(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	if user.GithubToken.RefreshTokenExpiresIn.Before(time.Now()) {
		newToken, err := utils.RefreshUserAccessToken(user.GithubToken)
		log.Println(newToken)

		if err != nil {
			log.Println(err)
			return utils.RespondError(c, fiber.StatusInternalServerError, err)
		}

		user.GithubToken = newToken

		if err := database.DB.Save(&user).Error; err != nil {
			log.Println(err)
			return utils.RespondError(c, fiber.StatusInternalServerError, err)
		}
	}

	log.Println("token should have expired")

	repos, err := utils.GetUserRepos(user.GithubToken.AccessToken)

	if err != nil {
		log.Println(err)
		return utils.RespondError(c, fiber.StatusInternalServerError, errors.New("unable to get user's repositories"))
	}

	return utils.RespondSuccess(c, repos, "Successfully retried user repos")
}

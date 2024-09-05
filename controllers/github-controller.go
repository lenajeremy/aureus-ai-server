package controllers

import (
	log "code-review/logger"
	"code-review/utils"
	"errors"
	// "time"

	"github.com/gofiber/fiber/v2"
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
	log.Logger.Infof("user: %v \n err: %v", user, err)
	log.Logger.Infoln(user.GithubToken)

	if err != nil {
		log.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	// if user.GithubToken.RefreshTokenExpiresIn.Before(time.Now()) {

	// 	log.Logger.Info("token should have expired")

	// 	err := utils.RefreshUserAccessToken(user.GithubToken)
	// 	log.Logger.Infof("new token: %v", user.GithubToken)
	// 	log.Logger.Info(user.GithubToken)

	// 	if err != nil {
	// 		log.Logger.Info(err)
	// 		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	// 	}

	// 	if err := database.DB.Save(&user).Error; err != nil {
	// 		log.Logger.Info(err)
	// 		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	// 	}
	// }

	repos, err := utils.GetUserRepos(user.GithubToken.AccessToken)

	if err != nil {
		log.Logger.Info(err)
		return utils.RespondError(c, fiber.StatusInternalServerError, errors.New("unable to get user's repositories"))
	}

	return utils.RespondSuccess(c, repos, "Successfully retried user repos")
}

package controllers

import (
	"code-review/config"
	"code-review/database"
	"code-review/models"
	"code-review/utils"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"net/url"
	"time"
)

func InitiateGHLogin(c *fiber.Ctx) error {
	sessId := c.Query("identifier")
	onSuccessCallback := c.Query("onSuccessCallback")
	onErrorCallback := c.Query("onErrorCallback")

	var successCallback, errorCallback url.URL

	if sessId == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("invalid session identifier"))
	}

	// verify the validity of success and error callback urls
	if sCallback, err := url.Parse(onSuccessCallback); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, err)
	} else {
		successCallback = *sCallback
	}
	if errCallback, err := url.Parse(onErrorCallback); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, err)
	} else {
		errorCallback = *errCallback
	}

	// generate state hash and save to db
	ghLoginStateHash := utils.Encrypt(sessId)
	err := database.DB.Create(&models.LoginInitSession{
		Hash:       ghLoginStateHash,
		Identifier: sessId,
		OnSuccess:  successCallback.String(),
		OnError:    errorCallback.String(),
	}).Error

	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	ghClientID := config.GetEnv("GITHUB_CLIENT_ID")

	redirectUri := c.BaseURL() + "/auth/github/callback"

	ghOAuthURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&state=%s&redirect_uri=%s",
		ghClientID, ghLoginStateHash, redirectUri,
	)

	return utils.RespondSuccess(c, ghOAuthURL, "Github OAuth URL")
}

func GHLoginCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	ghStateHash := c.Query("state")

	ghStateIdentifier := utils.Decrypt(ghStateHash)

	var loginSessState models.LoginInitSession

	err := database.DB.First(&loginSessState, "identifier = ? and hash = ?", ghStateIdentifier, ghStateHash).Error
	// delete the login session after the user has logged in successfully
	defer func() {
		if err := database.DB.Delete(&loginSessState).Error; err != nil {
			log.Println(err.Error())
		}
	}()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.RespondError(c, fiber.StatusBadRequest, err)
		} else {
			return utils.RespondError(c, fiber.StatusInternalServerError, err)
		}
	}

	// get token from code
	token, err := utils.ExchangeCodeForToken(code)
	if err != nil {
		return c.Redirect(fmt.Sprintf("%ss?error=%s", loginSessState.OnError, err.Error()))
	}

	// get user details from using token and
	user, err := utils.GetUserDetails(token.AccessToken, "")
	if err != nil {
		return c.Redirect(fmt.Sprintf("%ss?error=%s", loginSessState.OnError, err.Error()))
	}

	// save user to database
	dbUser := models.User{
		Name:          *user.Name,
		Email:         *user.Email,
		EmailVerified: time.Now(),
		Image:         *user.AvatarURL,
		GithubToken:   token,
	}

	if err := database.DB.FirstOrCreate(&dbUser).Error; err != nil {
		log.Println(err.Error())
		return c.Redirect(fmt.Sprintf("%ss?error=%s", loginSessState.OnError, err.Error()))
	}

	// generate jwt token
	jwtToken := utils.GenerateToken(dbUser.ID.String(), dbUser.Email)

	return c.Redirect(fmt.Sprintf("%s?token=%s", loginSessState.OnSuccess, jwtToken))
}

package auth

import (
	"code-review/config"
	"code-review/database"
	"code-review/github"
	"code-review/globals"
	"code-review/utils"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
)

func InitiateGHLogin(c *fiber.Ctx) error {
	sessId := c.Query("identifier")
	if sessId == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("invalid session identifier"))
	}

	// generate state hash and save to db
	ghLoginStateHash := utils.Encrypt(sessId)
	log.Println(ghLoginStateHash, sessId)
	err := database.DB.Create(&globals.LoginInitSession{Hash: ghLoginStateHash, Identifier: sessId}).Error

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

	var loginSessState globals.LoginInitSession

	err := database.DB.First(&loginSessState, "identifier = ? and hash = ?", ghStateIdentifier, ghStateHash).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.RespondError(c, fiber.StatusBadRequest, err)
		} else {
			return utils.RespondError(c, fiber.StatusInternalServerError, err)
		}
	}

	token, err := github.ExchangeCodeForToken(code)

	cookie := fiber.Cookie{
		Name:     "TOKEN",
		Value:    token,
		Secure:   false,
		HTTPOnly: false,
		SameSite: "none",
	}

	c.Cookie(&cookie)

	return c.Redirect(fmt.Sprintf("http://localhost:3000?token=%s", token))
}

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
	"time"
)

func InitiateGHLogin(c *fiber.Ctx) error {
	sessId := c.Query("identifier")
	if sessId == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("invalid session identifier"))
	}

	// generate state hash and save to db
	ghLoginStateHash := utils.Encrypt(sessId)
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
	token, err := github.ExchangeCodeForToken(code)
	if err != nil {
		return c.Redirect(fmt.Sprintf("http://localhost:3000/auth/error?error=%s", err.Error()))
	}

	// get user details from using token and
	user, err := github.GetUserDetails(token.AccessToken, "")
	if err != nil {
		return c.Redirect(fmt.Sprintf("http://localhost:3000/auth/error?error=%s", err.Error()))
	}

	// save user to database
	dbUser := User{
		Name:          *user.Name,
		Email:         *user.Email,
		EmailVerified: time.Now(),
		Image:         *user.AvatarURL,
		GithubToken:   token,
	}

	if err := database.DB.FirstOrCreate(&dbUser).Error; err != nil {
		log.Println(err.Error())
		return c.Redirect(fmt.Sprintf("http://localhost:3000/auth/error?error=%s", err.Error()))
	}

	log.Println(dbUser.ID)
	// generate jwt token
	jwtToken := GenerateToken(dbUser.ID.String(), dbUser.Email)

	return c.Redirect(fmt.Sprintf("http://localhost:3000?token=%s", jwtToken))
}

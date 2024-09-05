package controllers

import (
	"code-review/config"
	"code-review/database"
	"code-review/emails"
	"code-review/logger"
	"code-review/models"
	"code-review/utils"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

type emailsignup struct {
	FullName     string `json:"fullName"`
	Email        string `json:"email"`
	OnSuccessURL string `json:"onSuccessURL"`
	OnErrorURL   string `json:"onErrorURL"`
}

type emailsignin struct {
	Email        string `json:"email"`
	OnSuccessURL string `json:"onSuccessURL"`
	OnErrorURL   string `json:"onErrorURL"`
}

// EmailSignUp creates a new user and sends an email verification token
func EmailSignUp(c *fiber.Ctx) error {
	var body emailsignup

	if err := c.BodyParser(&body); err != nil {
		logger.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(body.Email) == "" || strings.TrimSpace(body.FullName) == "" || strings.TrimSpace(body.OnSuccessURL) == "" || strings.TrimSpace(body.OnErrorURL) == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("invalid email or full name"))
	}

	user := models.User{
		Name:  body.FullName,
		Email: body.Email,
	}

	err := database.DB.Create(&user).Error
	if err != nil {
		logger.Logger.Error("Error creating user: " + err.Error())
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value") {
			return utils.RespondError(c, fiber.StatusBadRequest, errors.New("user already exists"))
		} else {
			return utils.RespondError(c, fiber.StatusInternalServerError, err)
		}
	}

	emailVerificationToken := models.VerificationToken{
		Token:        emails.GenerateVerificationToken(body.Email),
		Email:        user.Email,
		Expires:      time.Now().Add(time.Minute * 10),
		Type:         models.VerificationTypeSignUp,
		OnSuccessURL: body.OnSuccessURL,
		OnErrorURL:   body.OnErrorURL,
	}

	logger.Logger.Infoln(emailVerificationToken, "creating email verification token")
	err = database.DB.Create(&emailVerificationToken).Error
	if err != nil {
		logger.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	redirectUrl := fmt.Sprintf("%s/auth/email/verify?token=%s&onSuccessURL=%s&onErrorURL=%s", c.BaseURL(), emailVerificationToken.Token, body.OnSuccessURL, body.OnErrorURL)
	logger.Logger.Infoln(redirectUrl, "redirect url")

	err = emails.SendEmailVerification(body.Email, body.FullName, redirectUrl)
	logger.Logger.Infoln(err, "error sending email")

	if err != nil {
		logger.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	return utils.RespondSuccess(c, nil, "Email verification sent")
}

func EmailSignIn(c *fiber.Ctx) error {

	var body emailsignin

	if err := c.BodyParser(&body); err != nil {
		logger.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(body.Email) == "" || strings.TrimSpace(body.OnSuccessURL) == "" || strings.TrimSpace(body.OnErrorURL) == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("invalid email or callback urls"))
	}

	var emailVerificationToken models.VerificationToken

	emailVerificationToken.Token = emails.GenerateVerificationToken(body.Email)
	emailVerificationToken.Email = body.Email
	emailVerificationToken.Expires = time.Now().Add(time.Minute * 10)
	emailVerificationToken.Type = models.VerificationTypeSignIn

	var user models.User
	err := database.DB.First(&user, "email = ?", body.Email).Error

	if err != nil {
		logger.Logger.Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.RespondError(c, fiber.StatusBadRequest, err)
		} else {
			return utils.RespondError(c, fiber.StatusInternalServerError, err)
		}
	}

	err = database.DB.Create(&emailVerificationToken).Error
	if err != nil {
		logger.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	redirectUrl := fmt.Sprintf("%s/auth/email/verify?token=%s&onSuccessURL=%s&onErrorURL=%s", c.BaseURL(), emailVerificationToken.Token, body.OnSuccessURL, body.OnErrorURL)

	err = emails.SendEmailVerification(body.Email, user.Name, redirectUrl)

	if err != nil {
		logger.Logger.Error(err.Error())
		return utils.RespondError(c, fiber.StatusInternalServerError, err)
	}

	return utils.RespondSuccess(c, nil, "Email verification sent")
}

func VerifyEmail(c *fiber.Ctx) error {
	log.Println("Verify email called")
	token := c.Query("token")
	onSuccessURL := c.Query("onSuccessURL")
	onErrorURL := c.Query("onErrorURL")

	logger.Logger.Infoln(token, onSuccessURL, onErrorURL)

	if token == "" || onSuccessURL == "" || onErrorURL == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("invalid token or callback urls"))
	}

	var verificationToken models.VerificationToken

	err := database.DB.First(&verificationToken, "token = ?", token).Error
	if err != nil {
		logger.Logger.Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect(fmt.Sprintf("%s?error=%s", onErrorURL, "invalid token"))
		} else {
			return c.Redirect(fmt.Sprintf("%s?error=%s", onErrorURL, "internal server error"))
		}
	}

	defer func() {
		if err := database.DB.Delete(&verificationToken).Error; err != nil {
			logger.Logger.Error(err.Error())
		}
	}()

	if verificationToken.Expires.Before(time.Now()) {
		logger.Logger.Error("token expired")
		return c.Redirect(fmt.Sprintf("%s?error=%s", onErrorURL, "token expired"))
	}

	var user models.User

	err = database.DB.First(&user, "email = ?", verificationToken.Email).Error
	if err != nil {
		logger.Logger.Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect(fmt.Sprintf("%s?error=%s", onErrorURL, "invalid token"))
		} else {
			return c.Redirect(fmt.Sprintf("%s?error=%s", onErrorURL, "internal server error"))
		}
	}

	if user.EmailVerified == nil {
		verifiedAt := time.Now()
		user.EmailVerified = &verifiedAt

		err = database.DB.Save(&user).Error

		if err != nil {
			logger.Logger.Error(err.Error())
			return c.Redirect(fmt.Sprintf("%s?error=%s", onErrorURL, "internal server error"))
		}
	}

	jwtToken := utils.GenerateJWTToken(user.ID.String(), user.Email)

	return c.Redirect(fmt.Sprintf("%s?token=%s", onSuccessURL, jwtToken))
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
		return c.Redirect(fmt.Sprintf("%s?error=%s", loginSessState.OnError, err.Error()))
	}

	// get user details from using token and
	user, err := utils.GetUserDetails(token.AccessToken, "")
	if err != nil {
		return c.Redirect(fmt.Sprintf("%s?error=%s", loginSessState.OnError, err.Error()))
	}

	var dbUser models.User

	if err := database.DB.First(&models.User{}, "email = ?", user.Email).Scan(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create new user
			verifiedAt := time.Now()
			dbUser = models.User{
				Name:          *user.Name,
				Email:         *user.Email,
				EmailVerified: &verifiedAt,
				Image:         *user.AvatarURL,
				// GithubToken:   &token,
			}

			if err := database.DB.Create(&dbUser).Error; err != nil {
				return c.Redirect(fmt.Sprintf("%s?error=%s", loginSessState.OnError, err.Error()))
			}
		} else {
			// update existing user github token
			// dbUser.GithubToken = &token
			if err := database.DB.Save(&dbUser).Error; err != nil {
				return c.Redirect(fmt.Sprintf("%s?error=%s", loginSessState.OnError, err.Error()))
			}
		}
	}

	// generate jwt token
	// jwtToken := utils.GenerateToken(dbUser.ID.String(), dbUser.Email)

	// ERROR: error
	return c.Redirect(fmt.Sprintf("%s?token=%s", loginSessState.OnSuccess, dbUser.ID.String()))
}

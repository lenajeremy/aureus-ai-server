package utils

import (
	"code-review/config"
	"code-review/database"
	"code-review/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"time"
	"errors"
	"fmt"
)

func GenerateJWTToken(id, email string) string {
	j := jwt.New(jwt.SigningMethodHS256)
	claims := j.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token, err := j.SignedString([]byte(config.GetEnv("JWT_SECRET")))

	if err != nil {
		log.Println(err.Error())
	}

	return token
}

func GetUserFromContext(c *fiber.Ctx) (*models.User, error) {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token in context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	var user models.User
	err = database.DB.Preload("GithubToken").First(&user, uid).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from database: %w", err)
	}

	return &user, nil
}

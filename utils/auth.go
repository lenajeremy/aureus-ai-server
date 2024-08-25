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
)

func GenerateToken(id, email string) string {
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

func GetUserFromContext(c *fiber.Ctx) (user models.User, err error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid, _ := uuid.Parse(claims["id"].(string))
	err = database.DB.Preload("GithubToken").First(&user, uid).Error

	return
}

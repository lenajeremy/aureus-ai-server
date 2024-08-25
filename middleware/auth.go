package middleware

import (
	"code-review/config"
	"code-review/utils"
	"errors"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.GetEnv("JWT_SECRET")),
		},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return utils.RespondError(c, fiber.StatusBadRequest, errors.New("missing or malformed JWT"))
	}
	log.Println(err)
	return utils.RespondError(c, fiber.StatusUnauthorized, errors.New("invalid or expired JWT"))
}

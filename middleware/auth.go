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
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err.Error())
			return utils.RespondError(c, fiber.StatusUnauthorized, errors.New("unauthorized, please login before you can access this endpoint"))
		},
	})
}

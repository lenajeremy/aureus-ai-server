package utils

import (
	"code-review/config"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

func GenerateToken(id, email string) string {
	j := jwt.New(jwt.SigningMethodHS256)
	claims := j.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24)

	token, err := j.SignedString([]byte(config.GetEnv("JWT_SECRET")))

	if err != nil {
		log.Println(err.Error())
	}

	return token
}

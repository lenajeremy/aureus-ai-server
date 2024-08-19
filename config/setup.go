package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err.Error())
	}
}

func GetEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	} else {
		return ""
	}
}

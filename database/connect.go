package database

import (
	cfg "code-review/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DBConnectConfig struct {
	MakeMigrations  bool
	MigrationModels []any
}

func ConnectToDB(config *DBConnectConfig) {
	var (
		DbHost = cfg.GetEnv("DB_HOST")
		DbPass = cfg.GetEnv("DB_PASS")
		DbName = cfg.GetEnv("DB_NAME")
		DbPort = cfg.GetEnv("DB_PORT")
		DbUser = cfg.GetEnv("DB_USER")
	)

	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", DbHost, DbUser, DbPass, DbPort, DbName)
	pg := postgres.Open(dsn)

	DB, err = gorm.Open(pg)

	if err != nil {
		log.Panicf("Failed to connect to database. Reason: %s", err.Error())
	} else {
		log.Println("Successfully connected to database")
	}

	if config != nil && config.MakeMigrations {

		err := DB.AutoMigrate(config.MigrationModels...)
		if err != nil {
			log.Panicf(err.Error())
		}
	}
}

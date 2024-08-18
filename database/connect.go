package database

import (
	"code-review/auth"
	cfg "code-review/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DBConnectConfig struct {
	MakeMigrations bool
}

func ConnectToDB(config ...DBConnectConfig) {

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

	db, err = gorm.Open(pg)

	if err != nil {
		log.Fatalf("Failed to connect to database. Reason: %s", err.Error())
	} else {
		log.Println("Successfully connected to database")
	}

	var dbConCfg DBConnectConfig

	if len(config) == 0 {
		dbConCfg = DBConnectConfig{
			MakeMigrations: true,
		}
	} else {
		dbConCfg = config[0]
	}

	if dbConCfg.MakeMigrations {
		err = db.AutoMigrate(
			&auth.Account{}, auth.Session{}, auth.User{}, auth.VerificationToken{},
		)
		if err != nil {
			log.Fatalf("Failed to make migrations. Reason: %s", err.Error())
		} else {
			log.Println("DB migrations completed")
		}
	}

}

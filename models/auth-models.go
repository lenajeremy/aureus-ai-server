package models

import (
	"time"
)

type User struct {
	BaseModel
	Name          string    `gorm:"type:VARCHAR(255)"`
	Email         string    `gorm:"type:VARCHAR(255);unique"`
	EmailVerified time.Time `gorm:"column:emailVerified"`
	Image         string    `gorm:"type:TEXT"`
	GithubToken   GHToken   `gorm:"foreignKey:user_id"`
}

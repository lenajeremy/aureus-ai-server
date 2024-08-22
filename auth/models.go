package auth

import (
	"code-review/github"
	"code-review/globals"
	"time"
)

type User struct {
	globals.BaseModel
	Name          string       `gorm:"type:VARCHAR(255)"`
	Email         string       `gorm:"type:VARCHAR(255);unique"`
	EmailVerified time.Time    `gorm:"column:emailVerified"`
	Image         string       `gorm:"type:TEXT"`
	GithubToken   github.Token `gorm:"foreignKey:user_id"`
}

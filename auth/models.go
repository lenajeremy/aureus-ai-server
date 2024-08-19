package auth

import (
	"code-review/globals"
	"time"
)

type User struct {
	globals.BaseModel
	Name          string    `gorm:"type:VARCHAR(255)"`
	Email         string    `gorm:"type:VARCHAR(255)"`
	EmailVerified time.Time `gorm:"column:emailVerified"`
	Image         string    `gorm:"type:TEXT"`
}

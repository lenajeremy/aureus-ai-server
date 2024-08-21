package github

import (
	"code-review/globals"
	"time"
)

type PRWebhookPayload struct {
	code          string
	affectedFiles []string
}

type Token struct {
	globals.BaseModel
	UserId                string    `gorm:"user_id"`
	AccessToken           string    `gorm:"access_token"`
	AccessTokenExpiresIn  time.Time `gorm:"expires_in;not null"`
	RefreshToken          string    `gorm:"refresh_token"`
	RefreshTokenExpiresIn time.Time `gorm:"refresh_token_expires_in;not nul"`
}

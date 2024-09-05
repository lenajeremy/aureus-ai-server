package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	BaseModel
	Name          string    `gorm:"type:VARCHAR(255);required"`
	Email         string    `gorm:"type:VARCHAR(255);unique;required"`
	EmailVerified *time.Time `gorm:"column:emailVerified"`
	Image         string    `gorm:"type:TEXT"`
	GithubToken   GHToken   `gorm:"foreignKey:ID"`
}

// Account is a user's account for a given authentication provider
type Account struct {
	BaseModel
	UserId            uuid.UUID    `gorm:"required"`
	User              User         `gorm:"foreignKey:UserId"`
	Provider          AuthProvider `gorm:"required"`
	ProviderAccountId string       `gorm:"required"`
	RefreshToken      string
	AccessToken       string
}

// AuthProvider indicates the method of authentication
type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGithub AuthProvider = "github"
	AuthProviderGoogle AuthProvider = "google"
)

// VerificationToken is used to verify a user's email address
type VerificationToken struct {
	BaseModel
	Token        string           `gorm:"required"`
	Expires      time.Time        `gorm:"required"`
	Email        string           `gorm:"required"`
	Type         VerificationType `gorm:"type:VARCHAR(10);required"`
	OnSuccessURL string           `gorm:"type:VARCHAR(255);required"`
	OnErrorURL   string           `gorm:"type:VARCHAR(255);required"`
}

type VerificationType string

const (
	VerificationTypeSignUp VerificationType = "signup"
	VerificationTypeSignIn VerificationType = "signin"
)

package auth

import (
	"code-review/globals"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type VerificationToken struct {
	globals.BaseModel
	Identifier string    `gorm:"primaryKey;type:TEXT;not null"`
	Expires    time.Time `gorm:"type:TIMESTAMPTZ;not null"`
	Token      string    `gorm:"primaryKey;type:TEXT;not null"`
}

type Account struct {
	globals.BaseModel
	UserID            int    `gorm:"column:userId;not null"`
	Type              string `gorm:"type:VARCHAR(255);not null"`
	Provider          string `gorm:"type:VARCHAR(255);not null"`
	ProviderAccountID string `gorm:"column:providerAccountId;type:VARCHAR(255);not null"`
	RefreshToken      string `gorm:"type:TEXT"`
	AccessToken       string `gorm:"type:TEXT"`
	ExpiresAt         int64  `gorm:"type:BIGINT"`
	IDToken           string `gorm:"column:id_token;type:TEXT"`
	Scope             string `gorm:"type:TEXT"`
	SessionState      string `gorm:"column:session_state;type:TEXT"`
	TokenType         string `gorm:"column:token_type;type:TEXT"`
}

type Session struct {
	globals.BaseModel
	UserID       int       `gorm:"column:userId;not null"`
	Expires      time.Time `gorm:"type:TIMESTAMPTZ;not null"`
	SessionToken string    `gorm:"column:sessionToken;type:VARCHAR(255);not null"`
}

type User struct {
	globals.BaseModel
	Name          string    `gorm:"type:VARCHAR(255)"`
	Email         string    `gorm:"type:VARCHAR(255)"`
	EmailVerified time.Time `gorm:"column:emailVerified;type:TIMESTAMPTZ"`
	Image         string    `gorm:"type:TEXT"`
}

func main() {
	// Replace with your actual database connection string
	dsn := "host=localhost user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&VerificationToken{}, &Account{}, &Session{}, &User{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Use the database connection (db) as needed
}

package globals

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;not null;type:uuid;unique;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"createdAt" gorm:"created_at;autoCreateTime:nano;not null;default:now()"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"updated_at;autoUpdateTime:nano;not null;default:now()"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"deleted_at;index"`
}

func (t *BaseModel) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}

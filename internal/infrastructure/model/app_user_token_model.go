package infrastructuremodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppUserToken represents the data structure for the app_user_token table.
type AppUserToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AppUserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_app_user_token_type" json:"appUserId"`
	TokenType string    `gorm:"type:varchar(25);default:'refresh';uniqueIndex:idx_app_user_token_type" json:"tokenType"`
	TokenUser string    `gorm:"type:varchar(200);not null" json:"tokenUser"`
	ExpireAt  time.Time `gorm:"type:timestamp;not null" json:"expireAt"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName explicitly overrides the table name used by GORM.
func (AppUserToken) TableName() string {
	return "app_user_token"
}

// BeforeCreate is a GORM hook that automatically generates a UUID if it is empty.
func (aut *AppUserToken) BeforeCreate(tx *gorm.DB) (err error) {
	if aut.ID == uuid.Nil {
		aut.ID = uuid.New()
	}
	return nil
}

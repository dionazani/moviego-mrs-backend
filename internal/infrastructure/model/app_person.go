package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppPerson represents the data structure for the "app_person" table in PostgreSQL.
// This table represents the core identity of the application user.
type AppPerson struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	SignUpFrom  *string    `gorm:"type:char(3)" json:"sign_up_from"` // web (WEB) or mobile (MBL)
	SignUpAt    *time.Time `gorm:"type:timestamp" json:"sign_up_at"`
	Fullname    string     `gorm:"type:varchar(100);not null" json:"fullname"`
	Gender      string     `gorm:"type:char(1);not null" json:"gender"`
	Email       string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	MobilePhone string     `gorm:"type:varchar(25);uniqueIndex;not null" json:"mobile_phone"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// TableName explicitly returns the table name in the database.
// This prevents GORM from auto-pluralizing the table name (into app_people).
func (AppPerson) TableName() string {
	return "app_person"
}

// BeforeCreate is a GORM hook called before a record is inserted into the database.
// This method automatically generates a new UUID if the ID field is empty (uuid.Nil).
func (ap *AppPerson) BeforeCreate(tx *gorm.DB) (err error) {
	if ap.ID == uuid.Nil {
		ap.ID = uuid.New()
	}
	return nil
}

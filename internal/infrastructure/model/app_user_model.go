package infrastructuremodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppUser represents the login credentials and security details for an app user.
type AppUser struct {
	ID                     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AppPersonID            uuid.UUID `gorm:"type:uuid;unique;not null" json:"appPersonId"`
	MstUserRoleID          uuid.UUID `gorm:"type:uuid;not null" json:"mstUserRoleId"`
	AppPassword            string    `gorm:"type:varchar(300);not null" json:"-"`
	MustChangePassword     int       `gorm:"type:int;default:0" json:"mustChangePassword"`
	NextChangePasswordDate time.Time `gorm:"type:date" json:"nextChangePasswordDate"`
	IsLocked               int        `gorm:"type:int;default:0" json:"isLocked"`
	FailedAttemptCount     int        `gorm:"type:int;default:0" json:"failedAttemptCount"`
	LockoutUntil           *time.Time `gorm:"type:timestamp;default:null" json:"lockoutUntil"`
	CreatedAt              time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName explicitly maps the struct to the "app_user" table.
func (AppUser) TableName() string {
	return "app_user"
}

// BeforeCreate is a GORM hook that automatically generates a UUID if it is empty.
func (au *AppUser) BeforeCreate(tx *gorm.DB) (err error) {
	if au.ID == uuid.Nil {
		au.ID = uuid.New()
	}
	return nil
}

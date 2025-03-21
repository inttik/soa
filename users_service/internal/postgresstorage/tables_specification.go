package postgresstorage

import (
	"time"

	"github.com/google/uuid"
)

type userInfo struct {
	Id      uuid.UUID `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	Login   string    `gorm:"size:255;not null;unique"`
	Pass    string    `gorm:"size:255;not null"`
	Root    bool      `gorm:"not null"`
	RegDate time.Time `gorm:"not null;type:timestamp"`
}

func (userInfo) TableName() string {
	return "user_info"
}

type userProfile struct {
	Id         uuid.UUID `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	FirstName  string    `gorm:"size:255"`
	LastName   string    `gorm:"size:255"`
	ImageLink  string    `gorm:"size:255"`
	BirthDate  time.Time `gorm:"type:date"`
	Telephone  string    `gorm:"size:20"`
	Email      string    `gorm:"not null"`
	LastModify time.Time `gorm:"not null;type:timestamp"`
}

func (userProfile) TableName() string {
	return "user_profile"
}

package postgresstorage

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type postTable struct {
	PostId      uuid.UUID      `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	AuthorId    uuid.UUID      `gorm:""`
	Title       string         `gorm:"size:255;not null"`
	Content     string         `gorm:"not null"`
	IsPrivate   bool           `gorm:"not null"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	PublishDate time.Time      `gorm:"not null;type:timestamp"`
	LastModify  time.Time      `gorm:"not null;type:timestamp"`
}

func (postTable) TableName() string {
	return "post"
}

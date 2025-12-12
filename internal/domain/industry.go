package domain

import (
	"time"
)

type Industry struct {
	ID        int       `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(100);not null"`
	Slug      string    `json:"slug" db:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM
func (Industry) TableName() string {
	return "industries"
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	FirstName   string         `gorm:"type:varchar(100);not null"`
	LastName    string         `gorm:"type:varchar(100);not null"`
	Email       string         `gorm:"uniqueIndex;type:varchar(255);not null"`
	PhoneNumber string         `gorm:"type:varchar(20)"`
	Type        string         `gorm:"type:varchar(50);not null"` // patient|employee|admin
	Password    string         `gorm:"type:varchar(255);not null"`
	ClinicID    *int64         `gorm:"index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

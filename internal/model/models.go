package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	SessionID  string    `gorm:"uniqueIndex;not null"`
	Protocol   string    `gorm:"not null"`
	TargetHost string    `gorm:"not null"`
	TargetPort int       `gorm:"not null"`
	Username   string    `gorm:"not null"`
	Status     string    `gorm:"not null"`
	StartTime  time.Time `gorm:"not null"`
	EndTime    time.Time
	UserID     uint `gorm:"not null"`
}

type AuditLog struct {
	gorm.Model
	SessionID string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}

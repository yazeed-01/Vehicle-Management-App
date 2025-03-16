package models

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	ID           uint           `gorm:"primaryKey"`
	UserID       string         `gorm:"index"`
	SearchedUser string         `gorm:"index"`
	SearchTime   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	SearchCount  int            `gorm:"default:1"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

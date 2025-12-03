package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;size:50;not null"`
	Email        string `gorm:"size:100;not null"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// 关联关系
	Channels []Channel `gorm:"foreignKey:UserID"`
	// Topics   []Topic   `gorm:"foreignKey:UserID"`
	// Articles []Article `gorm:"foreignKey:UserID"`
	// Comments []Comment `gorm:"foreignKey:UserID"`
}

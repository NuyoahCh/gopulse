package models

import (
	"time"
)

// Channel 频道模型
type Channel struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	UserID      uint   `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// 关联关系
	User User `gorm:"foreignKey:UserID"`
	// Topics []Topic `gorm:"foreignKey:ChannelID"`
}

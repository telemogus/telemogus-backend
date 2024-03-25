package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ChatName  string `json:"chatName" binding:"required" gorm:"not null"`
	IsGroup   bool   `json:"isGroup" binding:"required" gorm:"not null"`
	CreatedAt time.Time
	Messages  []Message `gorm:"foreignKey:ChatID"`
	Members   []User    `gorm:"many2many:user_chats;"`
}

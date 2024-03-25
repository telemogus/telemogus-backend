package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" binding:"required" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"password" binding:"required" gorm:"not null"`
	CreatedAt    time.Time
	LastSeen     time.Time
	Chats        []Chat `gorm:"many2many:user_chats;"`
}

package models

import (
	"time"
)

type User struct {
	Base
	Username     string    `json:"username" binding:"required" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" binding:"required" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt"`
	LastSeen     time.Time `json:"lastSeen"`
	Chats        []Chat    `json:"chats,omitempty" gorm:"many2many:user_chats;"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ChatID    uint
	UserID    uint
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}

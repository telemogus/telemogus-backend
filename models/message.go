package models

import (
	"time"
)

type Message struct {
	Base
	ChatID    uint
	UserID    uint
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}

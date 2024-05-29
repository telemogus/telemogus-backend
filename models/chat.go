package models

type Chat struct {
	Base
	ChatName string `json:"chatName" binding:"required" gorm:"not null"`
	Members  []User `json:"members" gorm:"many2many:user_chats;"`
}

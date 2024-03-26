package handlers

import (
	"net/http"

	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/models"

	"github.com/gin-gonic/gin"
)

func CreateChat(c *gin.Context) {
	var input struct {
		ChatName string `json:"chatName"`
		IsGroup  bool   `json:"isGroup"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentUser models.User
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id is not included into the token"})
		return
	}

	db.DB.Where("id = ?", userId).Find(&currentUser)

	chat := models.Chat{ChatName: input.ChatName, IsGroup: input.IsGroup, Members: []models.User{currentUser}}
	db.DB.Create(&chat)

	c.JSON(http.StatusOK, gin.H{"message": "Chat created"})
}

func GetUserChats(c *gin.Context) {
	userId := uint(c.MustGet("userId").(float64))
	var chatIDs []uint
	db.DB.Table("user_chats").Where("user_id = ?", userId).Pluck("chat_id", &chatIDs)

	var chats []models.Chat
	db.DB.Where("id IN ?", chatIDs).Preload("Messages").Preload("Members").Find(&chats)

	c.JSON(http.StatusOK, chats)
}

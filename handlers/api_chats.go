package handlers

import (
	"net/http"
	"time"

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
	type ChatPreview struct {
		Id                  uint      `json:"id"`
		ChatName            string    `json:"chatName"`
		LastMessageContent  string    `json:"lastMessageContent"`
		LastMessageTime     time.Time `json:"lastMessageTime"`
		UnreadMessagesCount uint      `json:"unreadMessagesCount"`
	}

	userId := uint(c.MustGet("userId").(float64))
	var chatIDs []uint
	db.DB.Table("user_chats").Where("user_id = ?", userId).Pluck("chat_id", &chatIDs)

	var chats []models.Chat
	db.DB.Where("id IN ?", chatIDs).Preload("Members").Find(&chats)

	chatsNumber := len(chats)
	userChatPreviews := make([]ChatPreview, chatsNumber)

	for i := 0; i < chatsNumber; i++ {
		currentChatId := chats[i].Id

		userChatPreviews[i].Id = currentChatId
		userChatPreviews[i].ChatName = chats[i].ChatName

		var lastMessage models.Message
		db.DB.Model(&models.Message{}).Where("chat_id = ?", currentChatId).Order("created_at desc").Find(&lastMessage)

		userChatPreviews[i].LastMessageTime = lastMessage.CreatedAt
		userChatPreviews[i].LastMessageContent = lastMessage.Content

		var unreadMessagesCount uint
		db.DB.Model(&models.Message{}).Select("count(*) as count").Where("state", models.Received).Scan(&unreadMessagesCount)

		userChatPreviews[i].UnreadMessagesCount = unreadMessagesCount
	}

	c.JSON(http.StatusOK, userChatPreviews)
}

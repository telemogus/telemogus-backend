package handlers

import (
	"net/http"
	"strconv"
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

func GetChat(c *gin.Context) {
	chatId, strerr := strconv.ParseUint(c.Param("chatId"), 10, 32)

	if strerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": strerr.Error()})
		return
	}

	chat := models.Chat{}

	if err := db.DB.Model(models.Chat{}).Where("id = ?", chatId).Preload("Members").Find(&chat); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func GetUserChats(c *gin.Context) {
	type ChatPreview struct {
		Id                  uint      `json:"id"`
		ChatName            string    `json:"chatName"`
		LastMessageContent  string    `json:"lastMessage"`
		LastMessageTime     time.Time `json:"lastMessageTime"`
		UnreadMessagesCount uint      `json:"unreadMessagesCount"`
	}

	userId := uint(c.MustGet("userId").(float64))
	var chatIDs []uint
	if err := db.DB.Table("user_chats").Where("user_id = ?", userId).Pluck("chat_id", &chatIDs); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	var chats []models.Chat
	if err := db.DB.Where("id IN ?", chatIDs).Preload("Members").Find(&chats); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	chatsNumber := len(chats)
	userChatPreviews := make([]ChatPreview, chatsNumber)

	for i := 0; i < chatsNumber; i++ {
		currentChatId := chats[i].Id

		userChatPreviews[i].Id = currentChatId
		userChatPreviews[i].ChatName = chats[i].ChatName

		var lastMessage models.Message
		if err := db.DB.Model(&models.Message{}).Where("chat_id = ?", currentChatId).Order("updated_at desc").Find(&lastMessage); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		userChatPreviews[i].LastMessageTime = lastMessage.CreatedAt
		userChatPreviews[i].LastMessageContent = lastMessage.Content

		var unreadMessagesCount uint
		if err := db.DB.Model(&models.Message{}).Select("count(*) as count").Where("state", models.Received).Scan(&unreadMessagesCount); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		userChatPreviews[i].UnreadMessagesCount = unreadMessagesCount
	}

	c.JSON(http.StatusOK, userChatPreviews)
}

func AddChatMember(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
	}

	chatId, strerr := strconv.ParseUint(c.Param("chatId"), 10, 32)

	if err := c.ShouldBindJSON(&input); err != nil || strerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", input.Username).Find(&user); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	userChat := struct {
		ChatId uint64 `json:"chatId"`
		UserId uint   `json:"userId"`
	}{chatId, user.Id}

	if err := db.DB.Table("user_chats").Create(&userChat); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, userChat)
}

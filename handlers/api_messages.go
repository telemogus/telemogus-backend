package handlers

import (
	"net/http"
	"strconv"

	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/models"

	"github.com/gin-gonic/gin"
)

func GetChatMessages(c *gin.Context) {
	chatId := c.Param("chatId")

	var messages []models.Message
	db.DB.Where("chat_id = ?", chatId).Order("created_at DESC").Find(&messages)

	c.JSON(http.StatusOK, messages)
}

func CreateChatMessage(c *gin.Context) {
	var input struct {
		Content string `json:"content"`
	}

	chatId, strerr := strconv.ParseUint(c.Param("chatId"), 10, 32)

	if err := c.ShouldBindJSON(&input); err != nil || strerr != nil {
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

	message := models.Message{ChatId: uint(chatId), UserId: currentUser.Id, Content: input.Content, State: models.Received}
	db.DB.Create(&message)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

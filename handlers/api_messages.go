package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/models"

	"github.com/gin-gonic/gin"
)

// GetChatMessages godoc
// @Summary Get chat messages
// @Description Retrieve all messages for a specific chat by chat ID
// @Tags chat
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Success 200 {array} models.Message
// @Router /chat/{chatId}/messages [get]
func GetChatMessages(c *gin.Context) {
	chatId := c.Param("chatId")

	var messages []models.Message
	db.DB.Where("chat_id = ?", chatId).Order("created_at DESC").Find(&messages)

	c.JSON(http.StatusOK, messages)
}

// CreateChatMessage godoc
// @Summary Send a message in a chat
// @Description Create a new message in a specific chat by chat ID
// @Tags chat
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param input body string true "Message content"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /chat/{chatId}/messages [post]
func CreateChatMessage(c *gin.Context) {
	var input struct {
		Content string `json:"content"`
	}

	chatId, strerr := strconv.ParseUint(c.Param("chatId"), 10, 32)

	if err := c.ShouldBindJSON(&input); err != nil || strerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id is not included into the token"})
		return
	}

	db.DB.Where("id = ?", userId).Find(&user)

	message := models.Message{ChatId: uint(chatId), UserId: user.Id, Content: input.Content, State: models.Received}
	db.DB.Create(&message)

	user.LastSeen = time.Now()
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

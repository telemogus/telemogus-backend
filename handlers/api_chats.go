package handlers

import (
	"net/http"
	"time"

	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/models"
	"github.com/golang-jwt/jwt"

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

	// userId, exists := c.Get("userId")

	// if !exists {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "User id is not included into the token"})
	// 	return
	// }

	chat := models.Chat{ChatName: input.ChatName, IsGroup: input.IsGroup, CreatedAt: time.Now()}
	db.DB.Create(&chat)

	// userChats := models.UserChats{ChatId: chat.ID, UserId: uint(userId.(float64))}
	// db.DB.Create(&userChats)

	c.JSON(http.StatusOK, gin.H{"message": "Chat created"})
}

func GetUserChats(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userId := claims["userId"].(string)

	var chats []models.Chat
	db.DB.Where("user_id = ?", userId).Find(&chats)

	c.JSON(http.StatusOK, chats)
}

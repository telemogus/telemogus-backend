package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var userConnections = struct {
	sync.RWMutex
	connections map[uint]*websocket.Conn
}{connections: make(map[uint]*websocket.Conn)}

func WebsocketHandler(c *gin.Context) {
	userIDStr := c.Param("userId") // Assume user_id is passed as query param
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing user_id"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}

	var user models.User
	if err := db.DB.Where("id = ?", userID).Find(&user); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to find user with this id"})
		return
	}

	websocketConnection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to set websocket upgrade"})
		return
	}

	userConnections.Lock()
	userConnections.connections[uint(userID)] = websocketConnection
	userConnections.Unlock()

	defer func() {
		userConnections.Lock()
		delete(userConnections.connections, uint(userID))
		userConnections.Unlock()
		websocketConnection.Close()
	}()

	for {
		var message models.Message
		err := websocketConnection.ReadJSON(&message)
		if err != nil {
			break
		}
		message.CreatedAt = time.Now()

		broadcastMessage(&message)
	}
}

func broadcastMessage(message *models.Message) {
	userConnections.RLock()
	defer userConnections.RUnlock()
	for userID, conn := range userConnections.connections {
		if userID == message.Id {
			continue
		}
		err := conn.WriteJSON(message)
		if err != nil {
			conn.Close()
			delete(userConnections.connections, userID)
		}
	}
}

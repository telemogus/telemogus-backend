package main

import (
	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/handlers"
	"github.com/dgb35/telemogus_backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()

	r.POST("/signup", handlers.SignUp)
	r.POST("/login", handlers.Login)

	authorized := r.Group("/", utils.AuthRequired)
	{
		authorized.GET("/ws", handlers.WebsocketHandler)
		authorized.POST("/chats", handlers.CreateChat)
		authorized.GET("/chats", handlers.GetUserChats)
		authorized.GET("/chats/:chatId", handlers.GetChat)

		authorized.POST("/chats/:chatId/members", handlers.AddChatMember)

		authorized.GET("/chats/:chatId/messages", handlers.GetChatMessages)
		authorized.POST("/chats/:chatId/messages", handlers.CreateChatMessage)
	}

	r.Run() // Listen and serve on 0.0.0.0:8080
}

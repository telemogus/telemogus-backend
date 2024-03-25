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
		authorized.POST("/chat", handlers.CreateChat)
		authorized.GET("/chats", handlers.GetUserChats)
	}

	r.Run() // Listen and serve on 0.0.0.0:8080
}

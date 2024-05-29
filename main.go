package main

import (
	"github.com/dgb35/telemogus_backend/db"
	"github.com/dgb35/telemogus_backend/docs"
	"github.com/dgb35/telemogus_backend/handlers"
	"github.com/dgb35/telemogus_backend/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func main() {
	db.Init()

	r := gin.Default()

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

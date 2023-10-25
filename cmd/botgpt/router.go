package botgpt

import (
	"botgpt/internal/clients/telegram"
	"botgpt/internal/config"
	"botgpt/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	router := gin.New()
	router.Use(gin.Recovery())

	viper := config.GetConfig()
	limit := int64(viper.GetInt("upload.limit"))
	router.MaxMultipartMemory = limit << 20 // 4 MiB

	telegram.CreateOrGetTgClient()

	router.Static("/static", "./static")

	v1 := router.Group("/api/v1")
	v1.Use(middleware.HttpLoggerMiddleware)
	v1.Use(middleware.ExceptionMiddleware)

	v1.GET("/status", statusController.Status)

	chat := v1.Group("/chat")
	{
		chat.POST("/completions", chatController.CompleteChat)
	}

	webhook := v1.Group("/webhook")
	{
		webhook.POST("/line", webHookController.LineMessage)
		webhook.POST("/tg", webHookController.TgMessage)
		webhook.POST("/client", webHookController.ClientMessage)
		webhook.POST("/azure", webHookController.AzureNotification)
		webhook.POST("/sendtg", webHookController.SendToTelegram)

		webhook.POST("/update/prompt", webHookController.UpdatePrompt)
		webhook.POST("/update/group", webHookController.UpdateGroup)
	}

	// ChatGPT
	root := router.Group("/")
	root.Use(middleware.CheckHeaderMiddleware(), middleware.HttpLoggerMiddleware)

	conversationsGroup := root.Group("/conversations")
	{
		conversationsGroup.GET("", chatgptController.GetConversations)

		// PATCH is official method, POST is added for Java support
		conversationsGroup.PATCH("", chatgptController.ClearConversations)
		conversationsGroup.POST("", chatgptController.ClearConversations)
	}

	conversationGroup := root.Group("/conversation")
	{
		conversationGroup.POST("", chatgptController.CreateConversation)
		conversationGroup.POST("/gen_title/:id", chatgptController.GenerateTitle)
		conversationGroup.GET("/:id", chatgptController.GetConversation)

		// rename or delete conversation use a same API with different parameters
		conversationGroup.PATCH("/:id", chatgptController.UpdateConversation)
		conversationGroup.POST("/:id", chatgptController.UpdateConversation)

		conversationGroup.POST("/message_feedback", chatgptController.FeedbackMessage)
	}

	// misc
	root.GET("/models", chatgptController.GetModels).Use(middleware.CheckHeaderMiddleware())
	root.GET("/accounts/check", chatgptController.GetAccountCheck).Use(middleware.CheckHeaderMiddleware())

	return router

}

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

	v1 := router.Group("/api/v1")
	v1.Use(middleware.HttpLoggerMiddleware)
	v1.Use(middleware.ExceptionMiddleware)

	v1.GET("/status", statusController.Status)

	webhook := v1.Group("/webhook")
	{
		webhook.POST("/line", webHookController.LineMessage)
		webhook.POST("/tg", webHookController.TgMessage)
		webhook.POST("/client", webHookController.ClientMessage)

		webhook.POST("/update/prompt", webHookController.UpdatePrompt)
		webhook.POST("/update/group", webHookController.UpdateGroup)
	}

	return router

}

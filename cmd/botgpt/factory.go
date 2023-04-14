package botgpt

import (
	"botgpt/internal/ai"
	"botgpt/internal/controllers"
	"botgpt/internal/handler"
)

var statusController *controllers.StatusController
var webHookController *controllers.WebHookController

func RegisterFactory() {

	var aiProvider = ai.NewGpt3AiProvider()

	var telegramMessageHandler = handler.NewTelegramMessageHandler(aiProvider)

	statusController = controllers.NewStatusController()
	webHookController = controllers.NewWebHookController(telegramMessageHandler)

}

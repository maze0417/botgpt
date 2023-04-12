package botgpt

import (
	"botgpt/internal/ai"
	"botgpt/internal/clients/telegram"
	"botgpt/internal/controllers"
)

var statusController *controllers.StatusController
var webHookController *controllers.WebHookController

func registerFactory() {

	var aiProvider = ai.NewGpt3AiProvider()

	var aiSender = ai.NewGpt3AiSender(aiProvider)
	var telegramClient = telegram.NewTelegramClient(aiProvider, aiSender)

	statusController = controllers.NewStatusController()
	webHookController = controllers.NewWebHookController(telegramClient, aiSender)

}

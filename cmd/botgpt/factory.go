package botgpt

import (
	"botgpt/internal/ai"
	"botgpt/internal/clients/aws"
	"botgpt/internal/controllers"
	"botgpt/internal/controllers/chatgpt"
	"botgpt/internal/handler"
	"botgpt/internal/services"
)

var statusController *controllers.StatusController
var webHookController *controllers.WebHookController
var chatgptController *chatgpt.ChatgptController

func RegisterFactory() {

	var aiProvider = ai.NewOpenAiProvider()

	var messageHandler = handler.NewMessageHandler(aiProvider)

	var textToSpeech = aws.NewPolly()

	var appService = services.NewAppService(aiProvider, messageHandler)
	var lineService = services.NewLineService(aiProvider, messageHandler, textToSpeech)
	var telegramService = services.NewTelegramService(aiProvider, messageHandler, textToSpeech)

	statusController = controllers.NewStatusController()
	webHookController = controllers.NewWebHookController(telegramService, lineService, appService)
	chatgptController = chatgpt.NewChatgptController()
}

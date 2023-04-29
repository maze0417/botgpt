package interfaces

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/models"
)

type IAiProvider interface {
	GenerateImage(message string, userID string) (string, error)
	GenerateText(totalMessages []gpt3.Message, userID string) (error, string)
	Transcribe(filePath string) (string, error)
}

type IMessageHandler interface {
	Send(messageFrom string, isGroup bool, userID string, groupID string, replyMessage string) (error, *models.AiResponse)
}

type IMessageService interface {
	HandleIfText(input interface{})
	HandleIFVoice(input interface{})
}
type ITextToSpeech interface {
	TextToSpeech(text string, outputFile string, outputFormat string, lang string) (error, []byte)
	GetLangFromText(text string) string
}

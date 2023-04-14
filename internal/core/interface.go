package core

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
	HandleText(messageFrom string, isGroup bool, userID string, groupID string, replyMessage string) (error, *models.AiResponse)
	HandleVoice(fileID string, isGroup bool, userID string, groupID string, replyMessage string) (error, *models.AiResponse)
}

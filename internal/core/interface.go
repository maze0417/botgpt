package core

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/models"
)

type IAiProvider interface {
	GenerateImage(message string) (string, error)
	GenerateText(totalMessages []gpt3.Message, userID string) (error, string)
	Transcribe(filePath string) (string, error)
}

type IMessageService interface {
	Send(messageFrom string, isGroup bool, userID string, groupID string, replyMessage string) (error, *models.AiResponse)
}

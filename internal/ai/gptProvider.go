package ai

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/core"
)

type Gpt3Provider struct {
}

func (a Gpt3Provider) Transcribe(filePath string) (string, error) {
	return gpt3.Transcribe(filePath)
}

func NewGpt3AiProvider() core.IAiProvider {
	return &Gpt3Provider{}
}
func (a Gpt3Provider) GenerateImage(message string) (string, error) {
	return gpt3.GenerateImageGpt(message)
}

func (a Gpt3Provider) GenerateText(totalMessages []gpt3.Message, userID string) (error, string) {
	return gpt3.CompletionGpt(totalMessages, userID)
}

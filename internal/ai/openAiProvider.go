package ai

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/core"
	"botgpt/internal/utils"

	"log"
)

type OpenAiProvider struct {
}

func (a OpenAiProvider) Transcribe(filePath string) (string, error) {
	return gpt3.Transcribe(filePath)
}

func NewGpt3AiProvider() core.IAiProvider {
	return &OpenAiProvider{}
}
func (a OpenAiProvider) GenerateImage(message string, userID string) (string, error) {
	log.Printf("send image gpt with %v :: \n %s \n\n", userID, message)

	resp, err := gpt3.GenerateImageGpt(message)

	if err != nil {
		log.Printf("send %v to gpt got error message :: %s  \n \n ", userID, err)
		return resp, err
	}
	log.Printf("reply %v image message ::\n %s  \n\n ", userID, resp)
	return resp, err
}

func (a OpenAiProvider) GenerateText(totalMessages []gpt3.Message, userID string) (error, string) {
	log.Printf("send gpt with %v :: \n %s \n\n", userID, utils.Json(totalMessages, true))

	err, resp := gpt3.CompletionGpt(totalMessages, userID)

	if err != nil {
		log.Printf("send %v to gpt got error message :: %s  \n \n ", userID, err)
		return err, resp
	}
	log.Printf("reply %v  message ::\n %s  \n\n ", userID, resp)
	return err, resp
}

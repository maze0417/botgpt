package gpt3

import (
	"botgpt/internal/config"
	"botgpt/internal/utils"
	"context"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	gpt3Once sync.Once
	myOpenAI *openai.Client
)

const (
	System    = "system" // 目前system作用較小
	User      = "user"
	Assistant = "assistant"
	gpt3      = "gpt-3.5-turbo"
	gpt4      = "gpt-4"
)

func createGp3Client() *openai.Client {
	gpt3Once.Do(func() {
		c := config.GetConfig()
		token := c.GetString("openai.access_token")
		myOpenAI = openai.NewClient(token)
	})
	return myOpenAI
}

func CompletionGpt(totalMessages []Message, userID string, model string) (error, string) {

	ctx := context.Background()

	request := openai.ChatCompletionRequest{
		Model:            model,
		Messages:         convertToChatCompletionMessages(totalMessages),
		MaxTokens:        1024,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		Stream:           false,
		User:             userID,
	}
	log.Printf("request to gpt :: %v :: \n %s \n\n", userID, utils.Json(request, true))

	resp, err := createGp3Client().CreateChatCompletion(ctx, request)

	if err != nil {
		log.Errorln(err)
		return err, ""
	}
	return nil, resp.Choices[0].Message.Content
}
func CreateChatCompletionStream(totalMessages []Message, userID string) (error, *openai.ChatCompletionStream) {

	ctx := context.Background()

	resp, err := createGp3Client().CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:            gpt3,
		Messages:         convertToChatCompletionMessages(totalMessages),
		MaxTokens:        512,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		User:             userID,
	})

	if err != nil {
		log.Errorln(err)
		return err, nil
	}
	return nil, resp
}
func convertToChatCompletionMessages(totalMessages []Message) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, len(totalMessages))

	for _, msg := range totalMessages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return messages
}

func GenerateImageGpt(message string) (string, error) {
	ctx := context.Background()

	resp, err := createGp3Client().CreateImage(ctx, openai.ImageRequest{
		Prompt:  message,
		N:       1,
		Size:    "1024x1024",
		Quality: "standard",
		Model:   openai.CreateImageModelDallE3,
	})

	if err != nil {
		log.Errorln(err)
		return "", err
	}
	return resp.Data[0].URL, err

}
func Transcribe(filePath string) (string, error) {
	ctx := context.Background()

	resp, err := createGp3Client().CreateTranscription(ctx, openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filePath,
		Prompt:   "",
	})

	if err != nil {
		log.Errorln(err)
		return "", err
	}
	return resp.Text, nil

}

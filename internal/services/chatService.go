package services

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/interfaces"
	"botgpt/internal/models"
	"encoding/json"
	"io"
)

type ChatService struct {
	aiProvider     interfaces.IAiProvider
	messageHandler interfaces.IMessageHandler
}

func NewChatService(aiProvider interfaces.IAiProvider, messageHandler interfaces.IMessageHandler) *ChatService {
	return &ChatService{
		aiProvider:     aiProvider,
		messageHandler: messageHandler,
	}
}

func (l ChatService) HandleIfText(input interface{}) {
	update, ok := input.(*models.ChatMessage)
	if ok {
		l.HandleText(update)
	}
}

func (l ChatService) HandleIFVoice(input interface{}) {

}

func (l ChatService) HandleText(message *models.ChatMessage) {

	var totalMessages []gpt3.Message

	for _, item := range message.ChatRequestMessage.Messages {
		msg := gpt3.Message{
			Role:    item.Role,
			Content: item.Content,
		}
		totalMessages = append(totalMessages, msg)
	}

	err, resp := l.aiProvider.GenerateTextStream(totalMessages, message.ChatRequestMessage.User)

	defer resp.Close()
	if err != nil {
		return
	}

	for {
		res, err := resp.Recv()
		if err == io.EOF {
			message.Context.SSEvent("", "[DONE]")
			return
		}
		data, _ := json.Marshal(res)

		message.Context.SSEvent("", string(data))
		message.Context.Writer.Flush()
	}

}

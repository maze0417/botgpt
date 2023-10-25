package models

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type Message struct {
	MessageID   int              `json:"message_id"`
	Date        int              `json:"date"`
	Text        string           `json:"text"`
	UserID      string           `json:"user_id"`
	ContextChan chan interface{} `json:"-"`
}

type ChatMessage struct {
	ChatRequestMessage openai.ChatCompletionRequest
	Context            *gin.Context `json:"-"`
}

type SendTgRequest struct {
	Text   string `json:"text"`
	ChatID string `json:"chat_id"`
}

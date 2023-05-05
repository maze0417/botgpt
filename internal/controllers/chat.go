package controllers

import (
	"botgpt/internal/enum"
	"botgpt/internal/interfaces"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/internal/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
)

type ChatController struct {
	webChatService interfaces.IMessageService
}

func NewChatController(
	webChatService interfaces.IMessageService,
) *ChatController {
	return &ChatController{
		webChatService: webChatService,
	}
}

func (h ChatController) CompleteChat(c *gin.Context) {

	var message openai.ChatCompletionRequest
	err := c.BindJSON(&message)
	if err != nil {

		utils.SendResponse(http.StatusBadRequest, response.Failure(fmt.Sprintf("Error parse request body : %v", err), enum.FALIURE), c)
		return
	}
	if len(message.Messages) == 0 {
		utils.SendResponse(http.StatusOK, response.OKHasContent("no message"), c)
		return
	}

	input := models.ChatMessage{
		ChatRequestMessage: message,
		Context:            c,
	}

	h.webChatService.HandleIfText(&input)
	h.webChatService.HandleIFVoice(&input)

}

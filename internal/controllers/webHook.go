package controllers

import (
	"botgpt/internal/ai"
	"botgpt/internal/clients/line"
	"botgpt/internal/enum"
	"botgpt/internal/interfaces"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/internal/utils/response"

	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

type WebHookController struct {
	telegramService interfaces.IMessageService
	lineService     interfaces.IMessageService
	appService      interfaces.IMessageService
}

func NewWebHookController(
	telegramService interfaces.IMessageService,
	lineService interfaces.IMessageService,
	appService interfaces.IMessageService,
) *WebHookController {
	return &WebHookController{
		telegramService: telegramService,
		lineService:     lineService,
		appService:      appService,
	}
}

func (h WebHookController) ClientMessage(c *gin.Context) {

	var message models.Message
	err := c.BindJSON(&message)
	if err != nil {

		utils.SendResponse(http.StatusBadRequest, response.Failure(fmt.Sprintf("Error parse request body : %v", err), enum.FALIURE), c)
		return
	}
	if len(message.Text) == 0 {
		utils.SendResponse(http.StatusOK, response.OKHasContent("no message"), c)
		return
	}

	message.ContextChan = make(chan interface{}, 1)
	h.appService.HandleIfText(&message)
	h.appService.HandleIFVoice(&message)

	res := <-message.ContextChan
	utils.SendResponse(http.StatusOK, response.OKHasContent(res), c)
}
func (h WebHookController) LineMessage(c *gin.Context) {
	events, err := line.ParseRequest(c.Request)
	if err != nil {

		utils.SendResponse(http.StatusBadRequest, response.Failure(fmt.Sprintf("Error parse request body : %v", err), enum.FALIURE), c)

		return
	}
	h.lineService.HandleIfText(events)

	h.lineService.HandleIFVoice(events)

	utils.SendResponse(http.StatusOK, response.OK(), c)
}

func (h WebHookController) TgMessage(c *gin.Context) {

	var update tgbotapi.Update
	err := c.BindJSON(&update)
	if err != nil {

		utils.SendResponse(http.StatusOK, response.Failure(fmt.Sprintf("Error parse request body : %v", err), enum.FALIURE), c)
		return
	}
	if update.Message == nil {
		utils.SendResponse(http.StatusOK, response.OKHasContent("no update.message"), c)
		return
	}

	h.telegramService.HandleIfText(update)

	h.telegramService.HandleIFVoice(update)

	utils.SendResponse(http.StatusOK, response.OKHasContent("received Message"), c)
}

func (h WebHookController) UpdatePrompt(c *gin.Context) {
	prompt := c.PostForm("prompt")

	var ct = ai.CommandMap[ai.ChildrenTalker]
	ct.System = prompt

	ai.CommandMap[ai.ChildrenTalker] = ct

	result := response.Make(true, enum.SUCCESS, "success", ct)

	utils.SendResponse(http.StatusOK, result, c)
}
func (h WebHookController) UpdateGroup(c *gin.Context) {
	cmd := c.PostForm("cmd")
	groupID := c.PostForm("groupID")

	resp := ai.SetGroupMode(groupID, cmd)

	if resp == nil {
		result := response.Make(false, enum.FALIURE, "cmd or group not set", nil)
		utils.SendResponse(http.StatusOK, result, c)
		return
	}
	result := response.OKHasContent(resp)

	utils.SendResponse(http.StatusOK, result, c)
}

package controllers

import (
	"botgpt/internal/ai"
	"botgpt/internal/clients/line"
	"botgpt/internal/core"
	"botgpt/internal/enum"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/internal/utils/redisManager"
	"botgpt/internal/utils/response"

	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WebHookController struct {
	telegramHandler core.IMessageHandler
}

func NewWebHookController(telegramHandler core.IMessageHandler) *WebHookController {
	return &WebHookController{
		telegramHandler: telegramHandler,
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

	sendMessage := fmt.Sprintf("%s %s", ai.ChildrenTalker, message.Text)
	//userID := fmt.Sprintf("client:%s", message.UserID)
	userID := fmt.Sprintf("client:%s", "1234")
	groupID := "1234"

	redisKey := fmt.Sprintf("%s:%s", userID, sendMessage)

	resp, err := redisManager.GetAndCache(redisKey, func() (interface{}, error) {
		resp, err := h.telegramHandler.HandleText(sendMessage, false, userID, groupID, "")
		return err, resp
	})

	switch err := err.(type) {
	case nil:
		// no error occurred, continue with your logic
	case *utils.KnownError:
		errString := fmt.Sprintf("KnownError occurs: %v", err)
		log.Errorln(errString)
		utils.SendResponse(http.StatusInternalServerError, response.Failure(errString, enum.FALIURE), c)
		return
	default:
		// unknown error occurred, log the error
		errString := fmt.Sprintf("exception occurs: %v", err)
		log.Errorln(errString)
		utils.SendResponse(http.StatusInternalServerError, response.Failure(errString, enum.FALIURE), c)
		return
	}

	//gptResponse, ok := resp
	//if !ok {
	//	// 轉型失敗，因為底層實際類型不是 *ApiResponse
	//	// 你可以在這裡處理錯誤，比如拋出 panic 或返回錯誤訊息
	//	errString := fmt.Sprintf("cast type error occurs: from %v to AiResponse", resp)
	//	log.Errorln(errString)
	//	utils.SendResponse(http.StatusInternalServerError, response.Failure(errString, enum.FALIURE), c)
	//	return
	//}

	utils.SendResponse(http.StatusOK, response.OKHasContent(resp), c)
}
func (h WebHookController) LineMessage(c *gin.Context) {
	events, err := line.CreateLineClient().ParseRequest(c.Request)
	if err != nil {

		utils.SendResponse(http.StatusBadRequest, response.Failure(fmt.Sprintf("Error parse request body : %v", err), enum.FALIURE), c)

		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				isGroup := event.Source.Type == "group"

				userID := fmt.Sprintf("line:%s", event.Source.UserID)
				groupID := event.Source.UserID
				if isGroup {
					userID = fmt.Sprintf("%s:%s", userID, event.Source.GroupID)
					groupID = event.Source.GroupID
				}

				err, gptRes := h.telegramHandler.HandleText(message.Text, isGroup, userID, groupID, "")

				switch err := err.(type) {
				case nil:
					// no error occurred, continue with your logic
				case *utils.KnownError:
					// err is a KnownError, you can access its properties
					continue
				default:
					// unknown error occurred, log the error
					log.Errorln(err)
					_, _ = line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
					continue
				}

				if gptRes.IsImage {
					if _, err = line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewImageMessage(gptRes.Text, gptRes.Text)).Do(); err != nil {
						log.Print(err)
						continue
					}

				}

				if gptRes.IsText {
					if _, err = line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(gptRes.Text)).Do(); err != nil {
						log.Print(err)
						continue
					}

				}

			}
		}
	}

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

	message := update.Message

	if message == nil || len(message.Text) == 0 {
		return nil, nil
	}

	isReply := message.ReplyToMessage != nil
	var messageReply string
	if isReply {
		messageReply = message.ReplyToMessage.Text
	}

	var messageFrom string

	if len(update.Message.Text) > 0 {
		messageFrom = message.Text
	}

	if len(update.Message.Caption) > 0 {
		messageFrom = update.Message.Caption
	}

	userID := fmt.Sprintf("tg:%s:%v", message.From.UserName, message.Chat.ID)
	groupID := fmt.Sprintf("%v", update.Message.Chat.ID)

	_, _ = h.telegramHandler.HandleText(messageFrom, update.Message.Chat.IsGroup(), userID, groupID, messageReply)

	_, _ = h.telegramHandler.HandleVoice()

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

package services

import (
	"botgpt/internal/ai"
	"botgpt/internal/interfaces"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/internal/utils/redisManager"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type AppService struct {
	aiProvider     interfaces.IAiProvider
	messageHandler interfaces.IMessageHandler
}

func (l AppService) HandleIfText(input interface{}) {
	go func() {
		update, ok := input.(*models.Message)
		if ok {
			l.HandleText(update)
		}
	}()
}

func (l AppService) HandleIFVoice(input interface{}) {
	go func() {
		update, ok := input.(models.Message)
		if ok {
			l.HandleVoice(update)
		}
	}()
}

func (l AppService) HandleText(message *models.Message) {

	sendMessage := fmt.Sprintf("%s %s", ai.ChildrenTalker, message.Text)
	//userID := fmt.Sprintf("client:%s", message.UserID)
	userID := fmt.Sprintf("client:%s", "1234")
	groupID := "1234"

	redisKey := fmt.Sprintf("%s:%s", userID, sendMessage)

	resp, err := redisManager.GetAndCache(redisKey, func() (interface{}, error) {
		resp, err := l.messageHandler.Send(sendMessage, false, userID, groupID, "")
		return err, resp
	})

	switch err := err.(type) {
	case nil:
		// no error occurred, continue with your logic
	case *utils.KnownError:
		errString := fmt.Sprintf("KnownError occurs: %v", err)
		log.Errorln(errString)
		//utils.SendResponse(http.StatusInternalServerError, response.Failure(errString, enum.FALIURE), message.Context)
		message.ContextChan <- errString
	default:
		// unknown error occurred, log the error
		errString := fmt.Sprintf("exception occurs: %v", err)
		log.Errorln(errString)
		//utils.SendResponse(http.StatusInternalServerError, response.Failure(errString, enum.FALIURE), message.Context)
		message.ContextChan <- errString
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
	message.ContextChan <- resp

}

func (l AppService) HandleVoice(message models.Message) {

}

func NewAppService(aiProvider interfaces.IAiProvider, messageHandler interfaces.IMessageHandler) *AppService {
	return &AppService{
		aiProvider:     aiProvider,
		messageHandler: messageHandler,
	}
}

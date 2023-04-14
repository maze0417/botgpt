package services

import (
	"botgpt/internal/clients/line"
	"botgpt/internal/interfaces"
	"botgpt/internal/utils"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
)

type LineService struct {
	aiProvider     interfaces.IAiProvider
	messageHandler interfaces.IMessageHandler
}

func (l LineService) HandleIfText(input interface{}) {
	go func() {
		update, ok := input.([]*linebot.Event)
		if ok {
			l.HandleText(update)
		}
	}()
}

func (l LineService) HandleIFVoice(input interface{}) {
	go func() {
		update, ok := input.([]*linebot.Event)
		if ok {
			l.HandleVoice(update)
		}
	}()
}

func NewLineService(aiProvider interfaces.IAiProvider, aiSender interfaces.IMessageHandler) *LineService {
	return &LineService{
		aiProvider:     aiProvider,
		messageHandler: aiSender,
	}
}

func (l LineService) HandleText(events []*linebot.Event) {

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

				err, gptRes := l.messageHandler.Send(message.Text, isGroup, userID, groupID, "")

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
}

func (l LineService) HandleVoice(update []*linebot.Event) {

}

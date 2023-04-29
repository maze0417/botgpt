package services

import (
	"botgpt/internal/clients/aws"
	"botgpt/internal/clients/line"
	"botgpt/internal/interfaces"
	"botgpt/internal/utils"
	"fmt"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type LineService struct {
	aiProvider     interfaces.IAiProvider
	messageHandler interfaces.IMessageHandler
	textToSpeech   interfaces.ITextToSpeech
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

func NewLineService(aiProvider interfaces.IAiProvider, aiSender interfaces.IMessageHandler, textToSpeech interfaces.ITextToSpeech) *LineService {
	return &LineService{
		aiProvider:     aiProvider,
		messageHandler: aiSender,
		textToSpeech:   textToSpeech,
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
						log.Printf("%s reply error %s ", gptRes.Text, err)
						continue
					}

				}

				if gptRes.IsText {
					if _, err = line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(gptRes.Text)).Do(); err != nil {
						log.Printf("%s reply error %s ", gptRes.Text, err)

						continue
					}

				}

			}
		}
	}
}

func (l *LineService) HandleVoice(events []*linebot.Event) {

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.AudioMessage:
				l.handleAudioMessage(line.CreateLineClient(), event, message)
			}
		}
	}
}
func (l *LineService) handleAudioMessage(bot *linebot.Client, event *linebot.Event, message *linebot.AudioMessage) {
	// 下載語音檔案
	content, err := bot.GetMessageContent(message.ID).Do()
	if err != nil {
		log.Println("無法下載語音檔案:", err)
		return
	}
	defer content.Content.Close()

	// 將 content.Content 保存為本地 MP3 檔案
	format := "m4a"
	localFilename := fmt.Sprintf("%s%s.%s", utils.GetUploadDir(), uuid.New().String(), format)

	data, err := io.ReadAll(content.Content)
	if err != nil {
		log.Println("無法讀取語音檔案內容:", err)
		return
	}
	err = os.WriteFile(localFilename, data, 0644)
	if err != nil {
		log.Println("無法將語音檔案保存為本地 MP3 檔案:", err)
		return
	}
	defer os.Remove(localFilename)
	text, err := l.aiProvider.Transcribe(localFilename)
	if err != nil {
		return
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("語音轉文字結果(Transcriptions result): \n%s", text))

	userID := fmt.Sprintf("line:%s", event.Source.UserID)
	groupID := event.Source.UserID

	err, gptResponse := l.messageHandler.Send(text, false, userID, groupID, "")
	builder.WriteString("\n\nGPT：\n")

	switch err := err.(type) {
	case nil:
		// no error occurred, continue with your logic
		builder.WriteString(gptResponse.Text)

		line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(builder.String())).Do()

		lang := l.textToSpeech.GetLangFromText(gptResponse.Text)
		if len(lang) == 0 {
			lang = gptResponse.Text
		}

		outputFile := fmt.Sprintf("%s%s.%s", utils.GetUploadDir(), uuid.New().String(), polly.OutputFormatMp3)
		log.Printf("try convert text to voice %s \n", outputFile)

		err, data = l.textToSpeech.TextToSpeech(gptResponse.Text, outputFile, polly.OutputFormatMp3, lang)
		if err != nil {
			builder.WriteString(err.Error())
			log.Error(err)
			line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(builder.String())).Do()
			return
		}
		defer os.Remove(outputFile)
		s3Client := aws.NewS3()
		duration := 3000 // 語音檔案的播放持續時間，單位為毫秒

		line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage("sending...")).Do()

		audioFileURL, err := s3Client.Upload(outputFile, data)
		if err != nil {
			log.Println("無法上傳到S3:", err)
			return
		}

		_, err = line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewAudioMessage(audioFileURL, duration)).Do()

		//_, err = bot.ReplyMessage(
		//	event.ReplyToken,
		//	linebot.NewAudioMessage(audioFileURL, duration),
		//).Do()

		if err != nil {
			log.Println("無法發送語音檔案:", err)
		}

		return
	case *utils.KnownError:
		// err is a KnownError, you can access its properties
		log.Errorln(err)

		return
	default:
		// unknown error occurred, log the error
		log.Errorln(err)
		builder.WriteString(err.Error())
		line.CreateLineClient().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(builder.String())).Do()
		return
	}

}

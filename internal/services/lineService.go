package services

import (
	"botgpt/internal/clients/line"
	"botgpt/internal/interfaces"
	"botgpt/internal/utils"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
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

func (l LineService) HandleVoice(events []*linebot.Event) {

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.AudioMessage:
				handleAudioMessage(line.CreateLineClient(), event, message)
			}
		}
	}
}
func handleAudioMessage(bot *linebot.Client, event *linebot.Event, message *linebot.AudioMessage) {
	// 下載語音檔案
	content, err := bot.GetMessageContent(message.ID).Do()
	if err != nil {
		log.Println("無法下載語音檔案:", err)
		return
	}
	defer content.Content.Close()

	// 將 content.Content 保存為本地 MP3 檔案
	localFilename := "downloaded_audio.mp3"
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

	localAudioFilePath := "your_local_audio_file.mp3"
	//data, err := io.ReadAll(localAudioFilePath)
	//if err != nil {
	//	log.Println("無法讀取本地 MP3 檔案:", err)
	//	return
	//}

	// 上傳 MP3 檔案到 AWS S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("your_aws_region"),
	})
	if err != nil {
		log.Println("無法建立 AWS Session:", err)
		return
	}

	s3Client := s3.New(sess)

	bucketName := "your_s3_bucket_name"
	objectKey := filepath.Base(localAudioFilePath)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"), // 設置為公開可讀取
	})

	if err != nil {
		log.Println("無法上傳 MP3 檔案到 S3:", err)
		return
	}

	// 使用 S3 上傳後的 URL 發送語音訊息
	audioFileURL := "https://" + bucketName + ".s3.amazonaws.com/" + objectKey
	duration := 3000 // 語音檔案的播放持續時間，單位為毫秒

	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewAudioMessage(audioFileURL, duration),
	).Do()

	if err != nil {
		log.Println("無法發送語音檔案:", err)
	}
}

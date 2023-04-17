package services

import (
	"botgpt/internal/clients/aws"
	"botgpt/internal/clients/telegram"
	"botgpt/internal/interfaces"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/pkg/ffmpeg"
	"fmt"
	"github.com/aws/aws-sdk-go/service/polly"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

type TelegramService struct {
	aiProvider     interfaces.IAiProvider
	messageHandler interfaces.IMessageHandler
}

func NewTelegramService(aiProvider interfaces.IAiProvider, aiSender interfaces.IMessageHandler) *TelegramService {
	return &TelegramService{
		aiProvider:     aiProvider,
		messageHandler: aiSender,
	}
}
func (t TelegramService) HandleIfText(input interface{}) {
	go func() {
		update, ok := input.(tgbotapi.Update)
		if ok {
			_, _ = t.HandleText(update)
		}
	}()
}

func (t TelegramService) HandleIFVoice(input interface{}) {
	go func() {
		update, ok := input.(tgbotapi.Update)
		if ok {
			_, _ = t.HandleVoice(update)
		}
	}()
}

func (t TelegramService) HandleText(update tgbotapi.Update) (*models.AiResponse, error) {

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

	_ = telegram.SendBotAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	err, gptResponse := t.messageHandler.Send(messageFrom, update.Message.Chat.IsGroup(), userID, groupID, messageReply)

	switch err := err.(type) {
	case nil:
		// no error occurred, continue with your logic
	case *utils.KnownError:
		// err is a KnownError, you can access its properties
		log.Errorln(err)

		return nil, err
	default:
		// unknown error occurred, log the error
		log.Errorln(err)
		_, _ = telegram.ReplayToChat(message.Chat.ID, err.Error(), "", update.Message.MessageID)

		return nil, err
	}
	escapedMessage := gptResponse.Text
	if gptResponse.IsText && gptResponse.TgParseMode == telegram.MarkdownV2 {
		escapedMessage = telegram.EscapeMessage(escapedMessage)
	}
	if gptResponse.IsImage {
		gptResponse.TgParseMode = ""
	}

	_, err = telegram.ReplayToChat(message.Chat.ID, escapedMessage, gptResponse.TgParseMode, update.Message.MessageID)
	if err != nil {
		_, _ = telegram.ReplayToChat(message.Chat.ID, err.Error(), gptResponse.TgParseMode, update.Message.MessageID)
	}
	return gptResponse, err
}

func (t TelegramService) HandleVoice(update tgbotapi.Update) (*models.AiResponse, error) {

	isNotVoice := update.Message == nil || update.Message.Voice == nil
	if isNotVoice {
		return nil, nil
	}
	message := update.Message

	groupID := fmt.Sprintf("%v", update.Message.Chat.ID)
	userID := fmt.Sprintf("tg:%s:%v", message.From.UserName, message.Chat.ID)

	inputOggName := fmt.Sprintf("%s%s-%d.ogg", utils.GetUploadDir(), groupID, update.Message.MessageID)
	format := polly.OutputFormatMp3
	outputName := fmt.Sprintf("%s%s-%d.%s", utils.GetUploadDir(), groupID, update.Message.MessageID, format)
	fileID := update.Message.Voice.FileID

	// retrieve the file using the Telegram file API
	_ = telegram.SendBotAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	voiceFileURL, err := telegram.CreateOrGetTgClient().GetFileDirectURL(fileID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	log.Printf("get voice url :: %s", voiceFileURL)
	voiceResponse, err := http.Get(voiceFileURL)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(voiceResponse.Body)

	voiceFile, err := os.Create(inputOggName)
	if err != nil {
		log.Printf("Error creating voice file: %s", err)
		return nil, err
	}
	defer func() {
		_ = voiceFile.Close()
		_ = os.Remove(inputOggName)
		_ = os.Remove(outputName)
	}()

	_, err = io.Copy(voiceFile, voiceResponse.Body)
	if err != nil {
		log.Printf("Error copying voice file: %s", err)
		return nil, err
	}

	log.Printf("Voice message downloaded to %s successfully , try  to convert mp3 %s \n", voiceFile.Name(), outputName)

	err = ffmpeg.ConvertOggToMp3(inputOggName, outputName)
	if err != nil {
		log.Println("ffmpeg Error converting OGG to MP3:", err)
		return nil, err
	}

	_ = telegram.SendBotAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	startPrompt := "語音辨識中(Processing your audio) ..."
	resMessage, _ := telegram.ReplayToChat(message.Chat.ID, startPrompt, "", update.Message.MessageID)

	text, err := t.aiProvider.Transcribe(outputName)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("語音轉文字結果(Transcriptions result): \n%s", text))

	_ = telegram.UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)

	err, gptResponse := t.messageHandler.Send(text, false, userID, groupID, "")
	builder.WriteString("\n\nGPT：\n")

	switch err := err.(type) {
	case nil:
		// no error occurred, continue with your logic
		builder.WriteString(gptResponse.Text)

		_ = telegram.UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)

		outputFile := fmt.Sprintf("%s%s.%s", utils.GetUploadDir(), uuid.New().String(), format)
		log.Printf("try convert text to voice %s \n", outputFile)

		err = aws.SynthesizeSpeech(gptResponse.Text, outputFile, format, gptResponse.CommandMode)
		if err != nil {
			builder.WriteString(err.Error())
			_ = telegram.UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)
			log.Error(err)

			return nil, err
		}
		err = telegram.SendVoice(outputFile)
		if err != nil {
			builder.WriteString(err.Error())
			_ = telegram.UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)
			log.Error(err)
			return nil, err
		}

		return gptResponse, nil
	case *utils.KnownError:
		// err is a KnownError, you can access its properties
		log.Errorln(err)

		return nil, err
	default:
		// unknown error occurred, log the error
		log.Errorln(err)
		builder.WriteString(err.Error())
		_ = telegram.UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)
		return nil, err
	}

}

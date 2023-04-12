package telegram

import (
	"botgpt/internal/config"
	"botgpt/internal/core"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/pkg/ffmpeg"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const MarkdownV2 = "MarkdownV2"

var (
	tgOnce   sync.Once
	tgClient *tgbotapi.BotAPI
)

type TelegramClient struct {
	aiProvider core.IAiProvider
	aiSender   core.IAiSender
}

func NewTelegramClient(aiProvider core.IAiProvider, aiSender core.IAiSender) *TelegramClient {
	return &TelegramClient{
		aiProvider: aiProvider,
		aiSender:   aiSender,
	}
}

var EscapeMarkdownString = []string{"_", "*",
	"[", "]",
	"(", ")",
	"~",
	"Pre", "pre",
	//"`",
	">",
	"#", "+", "-", "=", "|",
	"{", "}",
	".",
	"!",
}

func EscapeMessage(msg string) string {

	for _, v := range EscapeMarkdownString {
		msg = strings.ReplaceAll(msg, v, fmt.Sprintf("\\%s", v))
	}
	return msg
}

func CreateOrGetTgClient() *tgbotapi.BotAPI {
	tgOnce.Do(func() {
		c := config.GetConfig()
		token := c.GetString("tg.access_token")
		var err error

		tgClient, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Errorln(err)
		}

	})
	return tgClient

}
func ReplayToChat(chatID int64, msg string, tgParseMode string, replyMessageID int) (*tgbotapi.Message, error) {

	newMsg := tgbotapi.NewMessage(chatID, msg)
	if len(tgParseMode) > 0 {
		newMsg.ParseMode = tgParseMode
	}
	if replyMessageID > 0 {
		newMsg.ReplyToMessageID = replyMessageID
	}

	res, err := CreateOrGetTgClient().Send(newMsg)
	if err != nil {
		log.Errorf("Send telegram message error:: %v \n", err)
		return nil, err
	}
	return &res, nil
}
func SendBotAction(chatID int64, action string) error {

	newMsg := tgbotapi.NewChatAction(chatID, action)

	_, err := CreateOrGetTgClient().Send(newMsg)
	if err != nil {
		log.Errorf("Send telegram message error:: %v \n", err)
		return err
	}
	return nil
}
func (t TelegramClient) HandleText(update tgbotapi.Update) (*models.AiResponse, error) {

	message := update.Message

	if message == nil || len(message.Text) == 0 {
		return nil, nil
	}

	isReply := message.ReplyToMessage != nil

	var messageFrom string

	if len(update.Message.Text) > 0 {
		messageFrom = message.Text
	}

	if len(update.Message.Caption) > 0 {
		messageFrom = update.Message.Caption
	}

	if isReply {
		messageFrom = appendMessage(messageFrom, message.ReplyToMessage.Text)
	}

	userID := fmt.Sprintf("tg:%s:%v", message.From.UserName, message.Chat.ID)
	groupID := fmt.Sprintf("%v", update.Message.Chat.ID)

	_ = SendBotAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	err, gptResponse := t.aiSender.Send(messageFrom, update.Message.Chat.IsGroup(), userID, groupID)

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
		_, _ = ReplayToChat(message.Chat.ID, err.Error(), "", update.Message.MessageID)

		return nil, err
	}
	escapedMessage := gptResponse.Text
	if gptResponse.IsText && gptResponse.TgParseMode == MarkdownV2 {
		escapedMessage = EscapeMessage(escapedMessage)
	}
	if gptResponse.IsImage {
		gptResponse.TgParseMode = ""
	}

	_, err = ReplayToChat(message.Chat.ID, escapedMessage, gptResponse.TgParseMode, update.Message.MessageID)
	if err != nil {
		_, _ = ReplayToChat(message.Chat.ID, err.Error(), gptResponse.TgParseMode, update.Message.MessageID)
	}
	return gptResponse, err
}

func (t TelegramClient) HandleVoice(update tgbotapi.Update) (*models.AiResponse, error) {

	isNotVoice := update.Message == nil || update.Message.Voice == nil
	if isNotVoice {
		return nil, nil
	}
	message := update.Message

	groupID := fmt.Sprintf("%v", update.Message.Chat.ID)
	userID := fmt.Sprintf("tg:%s:%v", message.From.UserName, message.Chat.ID)

	inputOggName := fmt.Sprintf("v%s-%d.ogg", groupID, update.Message.MessageID)
	outputMp3Name := fmt.Sprintf("v%s-%d.mp3", groupID, update.Message.MessageID)
	fileID := update.Message.Voice.FileID

	// retrieve the file using the Telegram file API
	_ = SendBotAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	voiceFileURL, err := CreateOrGetTgClient().GetFileDirectURL(fileID)

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
		_ = os.Remove(outputMp3Name)
	}()

	_, err = io.Copy(voiceFile, voiceResponse.Body)
	if err != nil {
		log.Printf("Error copying voice file: %s", err)
		return nil, err
	}

	fmt.Printf("Voice message downloaded to %s successfully , try  to convert mp3 %s \n", voiceFile.Name(), outputMp3Name)

	err = ffmpeg.ConvertOggToMp3(inputOggName, outputMp3Name)
	if err != nil {
		fmt.Println("ffmpeg Error converting OGG to MP3:", err)
		return nil, err
	}

	_ = SendBotAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	startPrompt := "語音辨識中(Processing your audio) ..."
	resMessage, _ := ReplayToChat(message.Chat.ID, startPrompt, "", update.Message.MessageID)

	text, err := t.aiProvider.Transcribe(outputMp3Name)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("語音轉文字結果(Transcriptions result): \n%s", text))

	_ = UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)

	err, gptResponse := t.aiSender.Send(text, false, userID, groupID)
	builder.WriteString("\n\nGPT：\n")
	switch err := err.(type) {
	case nil:
		// no error occurred, continue with your logic
		builder.WriteString(gptResponse.Text)
		_ = UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)

		return gptResponse, nil
	case *utils.KnownError:
		// err is a KnownError, you can access its properties
		log.Errorln(err)

		return nil, err
	default:
		// unknown error occurred, log the error
		log.Errorln(err)
		builder.WriteString(err.Error())
		_ = UpdateMessage(message.Chat.ID, builder.String(), resMessage.MessageID)
		return nil, err
	}

}

func appendMessage(message string, append string) string {
	var builder strings.Builder
	builder.WriteString(message)
	builder.WriteRune('\n')

	builder.WriteString(append)
	builder.WriteRune('\n')

	return builder.String()
}

//
//func sendVoice() {
//
//	// Send the converted MP3 file back
//	mp3File, err := os.Open("output.mp3")
//	if err != nil {
//		fmt.Println("Error opening MP3 file:", err)
//		continue
//	}
//	defer mp3File.Close()
//
//	audioConfig := tgbotapi.NewAudioUpload(update.Message.Chat.ID, mp3File)
//	audioConfig.Title = "Converted Voice Message"
//	audioConfig.MimeType = "audio/mpeg"
//	audioConfig.FileID = "output.mp3"
//
//	_, err = bot.Send(audioConfig)
//	if err != nil {
//		fmt.Println("Error sending MP3 file:", err)
//	}
//}

func UpdateMessage(chatID int64, msg string, editMessageID int) error {

	newMsg := tgbotapi.NewEditMessageText(chatID, editMessageID, msg)

	_, err := CreateOrGetTgClient().Send(newMsg)
	if err != nil {
		log.Errorf("Send telegram message error:: %v \n", err)
		return err
	}
	return nil
}

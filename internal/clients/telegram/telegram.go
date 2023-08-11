package telegram

import (
	"botgpt/internal/config"
	"botgpt/pkg/ffmpeg"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
)

const MarkdownV2 = "MarkdownV2"

var (
	tgOnce   sync.Once
	tgClient *tgbotapi.BotAPI
)

var EscapeMarkdownString = []string{"_", "*",
	"[", "]",
	"(", ")",
	"~",
	"Pre", "pre", "PreCode", "precode",
	"<",
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
	//if len(tgParseMode) > 0 {
	//	newMsg.ParseMode = tgParseMode
	//}
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

func SendVoice(chatID int64, voicePath string) error {

	outputFile := strings.ReplaceAll(voicePath, ".mp3", ".ogg")

	err := ffmpeg.ConvertMp3ToOgg(voicePath, outputFile)
	if err != nil {
		return err
	}

	voice, err := os.Open(outputFile)
	if err != nil {
		return err
	}
	defer func(mp3File *os.File) {
		_ = mp3File.Close()
		_ = os.Remove(voicePath)
		_ = os.Remove(outputFile)
	}(voice)
	data, err := io.ReadAll(voice)

	audioConfig := tgbotapi.NewVoiceUpload(chatID, tgbotapi.FileBytes{Name: voice.Name(),
		Bytes: data})

	_, err = CreateOrGetTgClient().Send(audioConfig)
	if err != nil {
		fmt.Println("Error sending voice file:", err)
		return err
	}
	return nil
}

func UpdateMessage(chatID int64, msg string, editMessageID int) error {

	newMsg := tgbotapi.NewEditMessageText(chatID, editMessageID, msg)

	_, err := CreateOrGetTgClient().Send(newMsg)
	if err != nil {
		log.Errorf("Send telegram message error:: %v \n", err)
		return err
	}
	return nil
}

func SendMessage(chatID int64, msg string) (*tgbotapi.Message, error) {

	newMsg := tgbotapi.NewMessage(chatID, msg)

	res, err := CreateOrGetTgClient().Send(newMsg)
	if err != nil {
		log.Errorf("Send telegram message error:: %v \n", err)
		return nil, err
	}
	return &res, nil
}

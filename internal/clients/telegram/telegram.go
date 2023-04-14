package telegram

import (
	"botgpt/internal/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
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

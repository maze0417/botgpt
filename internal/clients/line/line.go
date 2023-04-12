package line

import (
	"botgpt/internal/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	once      sync.Once
	myLineBot *linebot.Client
)

func CreateLineClient() *linebot.Client {
	once.Do(func() {
		var err error

		c := config.GetConfig()
		token := c.GetString("linebot.access_token")
		secret := c.GetString("linebot.secret")

		myLineBot, err = linebot.New(
			secret,
			token,
		)
		if err != nil {
			log.Error(err)
		}
		fmt.Println("line client create success")
	})
	return myLineBot

}

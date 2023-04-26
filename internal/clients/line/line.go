package line

import (
	"botgpt/internal/config"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
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

// ParseRequest func
func ParseRequest(r *http.Request) ([]*linebot.Event, error) {

	if config.IsNotLocal() {
		return CreateLineClient().ParseRequest(r)
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err = json.Unmarshal(body, request); err != nil {
		return nil, err
	}
	return request.Events, nil

}

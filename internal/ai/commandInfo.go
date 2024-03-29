package ai

import (
	"botgpt/internal/clients/azure"
	"botgpt/internal/clients/telegram"
	"botgpt/internal/enum"
	"botgpt/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type CommandInfo struct {
	Cmd                string                       `json:"cmd"`
	PromptPrefixInject string                       `json:"prefix"`
	SystemPrompt       string                       `json:"prompt"`
	Usage              string                       `json:"usage"`
	TgParserMode       string                       `json:"tg_parser_mode"`
	PostExec           func(string) (string, error) `json:"-"`
	MaxHistoryLen      int                          `json:"max _history_len"`
	Lang               string                       `json:"-"`
	Alias              []string                     `json:"alias"`
}

const (
	ChatWithoutTag      = "/chatnotag@mazeaibot"
	Chat                = "/chat@mazeaibot"
	Image               = "/image@mazeaibot"
	CreateAzureWorkItem = "/cw@mazeaibot"
	ChildrenTalker      = "/ct@mazeaibot"
	EnToTw              = "/entw@mazeaibot"
	JpToTw              = "/jptw@mazeaibot"
	Help                = "/help@mazeaibot"
	Asker               = "/asker@mazeaibot"
)

const (
	UseGroupDefaultSysMsg = "你是一個有用AI助手，會嘗試回覆各種問題"
	//EnTwPrompt = "Translate into other language: If the language is Traditional Chinese, please translate it into English; If the language is English, please translate it into Traditional Chinese.Only print translated result without any additional information."
	//JpTwPrompt = "Translate into other language: If the language is Traditional Chinese, please translate it into Japanese; If the language is Japanese, please translate it into Traditional Chinese.Only print translated result without any additional information."
	EnTwPrompt = "translate belows into %s"
	JpTwPrompt = "translate belows into %s"
)

var (
	CommandMap = map[string]CommandInfo{
		Chat: {
			Cmd:           Chat,
			SystemPrompt:  UseGroupDefaultSysMsg, //use group as default
			PostExec:      nil,
			TgParserMode:  telegram.MarkdownV2,
			Usage:         "/chat {message}",
			MaxHistoryLen: 3,
			Alias:         []string{"@bot", "@mazeaibot", "/chat", "＠bot", "＠mazeaibot"},
		},
		Help: {
			Cmd:           Help,
			SystemPrompt:  "",
			TgParserMode:  "",
			Usage:         "/help show current mode and commands",
			MaxHistoryLen: 0,
			Alias:         []string{"/help"},
		},
		ChatWithoutTag: {
			Cmd:           ChatWithoutTag,
			SystemPrompt:  "",
			PostExec:      nil,
			TgParserMode:  telegram.MarkdownV2,
			Usage:         "/chatnotag chat without tag bot",
			MaxHistoryLen: 0,
			Alias:         []string{"/chatnotag"},
		},
		JpToTw: {
			Cmd:                JpToTw,
			SystemPrompt:       "",
			PromptPrefixInject: JpTwPrompt,
			TgParserMode:       "",
			Usage:              "/jptw@mazeaibot translate to ja-JP",
			MaxHistoryLen:      0,
			Lang:               enum.JaJP,
			Alias:              []string{"/jptw"},
		},
		EnToTw: {
			Cmd:                EnToTw,
			SystemPrompt:       "",
			PromptPrefixInject: EnTwPrompt,
			TgParserMode:       tgbotapi.ModeHTML,
			Usage:              "/entw@mazeaibot translate to en-US",
			MaxHistoryLen:      0,
			Lang:               enum.EnUS,
			Alias:              []string{"/entw"},
		},
		Image: {
			Cmd:          Image,
			SystemPrompt: UseGroupDefaultSysMsg,
			PostExec:     nil,
			//TgParserMode:  tgbotapi.ModeHTML,
			Usage:         "/image generate image by dall-e",
			MaxHistoryLen: 0,
			Alias:         []string{"@botimg", "draw", "畫", "/image", "＠botimg"},
		},

		CreateAzureWorkItem: {
			Cmd:           CreateAzureWorkItem,
			SystemPrompt:  "幫忙把把對話裡面的關於`標題` 跟 `內容`轉化成 json , 範例:幫我建立一張azure workitem 標題是`MyHeader` 內容是`Mycontent` 只需要輸出JSON 不需要其他內容: `[{\"op\": \"add\",\"path\": \"/fields/System.Title\",\"value\": \"MyHeader\"},{\"op\": \"add\",\"path\": \"/fields/System.Description\",\"value\": \"Mycontent\"}]` ",
			PostExec:      createAzureWorkItem,
			TgParserMode:  telegram.MarkdownV2,
			Usage:         "/cw@mazeaibot 開單 標題{title} 內容 {content} => 此指令可以用來開azure work item",
			MaxHistoryLen: 0,
			Alias:         []string{"/cw"},
		},
		ChildrenTalker: {
			Cmd:          ChildrenTalker,
			SystemPrompt: "你將扮演一個與小孩互動的虛擬人物:巧克力姐姐，不要稱呼自己是人工智能或AI。",
			//System:        "你是一位特立獨行、充滿愛心的人物，名為巧克力姐姐。你喜歡和小孩溝通，並且經常會將自己當作小孩的好朋友，開始一場充滿想像力的冒險之旅。你喜歡讓小孩開心，並且總是以最真誠的態度對待他們。你總是準備著一些小驚喜，例如準備了一份在小孩最喜愛的食物裡藏著神秘禮物的巧克力，這讓小孩嘴巴裡美味難忘，心底裡溫暖感激。你是一種慷慨大方、友善親切的存在，小孩總是對你充滿好感和信任。你的微笑和笑聲，能夠撫平所有痛苦和忧伤的心靈，這就是你所擁有的使命和魅力。",
			PostExec:           nil,
			TgParserMode:       telegram.MarkdownV2,
			Usage:              "/ct@mazeaibot {message} => ai baby",
			MaxHistoryLen:      3,
			PromptPrefixInject: "用小孩能夠聽懂的方式",
			Alias:              []string{"/aibaby"},
		},
		Asker: {
			Cmd:                Asker,
			SystemPrompt:       "",
			PostExec:           nil,
			TgParserMode:       telegram.MarkdownV2,
			Usage:              "/asker@mazeaibot 用是非題詢問",
			MaxHistoryLen:      0,
			PromptPrefixInject: "請你用是非題問我: \n",
			Alias:              []string{"/asker", "/ask", "@ask"},
		},
	}
)

func (c *CommandInfo) HaveHistoryMessage() bool {
	return c.MaxHistoryLen > 0
}

func GetCommandFromAlias(cmd string) *CommandInfo {

	for _, commandInfo := range CommandMap {
		if commandInfo.Cmd == cmd {
			return &commandInfo
		}
		if utils.Contains(commandInfo.Alias, cmd) {
			return &commandInfo
		}
	}

	return nil
}

func HasCommandPrefix(message string) bool {

	index := strings.Index(message, " ")

	// Substring up to the first space
	cmd := message
	if index < 0 {
		return false
	}

	if index >= 0 {
		cmd = message[0:index]
	}

	for _, commandInfo := range CommandMap {
		if commandInfo.Cmd == cmd {
			return true
		}
		if utils.Contains(commandInfo.Alias, cmd) {
			return true
		}
	}

	return false
}

func GetGroupModeOrCommandInfoByMessage(message string, groupID string) CommandInfo {

	groupMode := GetGroupMode(groupID)

	// Find the index position of the first space
	index := strings.Index(message, " ")

	// Substring up to the first space
	result := message
	if index >= 0 {
		result = message[0:index]
	}

	cmd := GetCommandFromAlias(result)

	if groupMode != nil {
		message = groupMode.CommandMode

		if cmd != nil && cmd.Cmd == Image {
			return CommandMap[Image]
		}
		return CommandMap[groupMode.CommandMode]
	}
	if cmd == nil {
		return CommandMap[Chat]
	}
	return CommandMap[cmd.Cmd]
}
func ReplaceCommandAsEmpty(msg string) string {

	for k, v := range CommandMap {
		for _, alias := range v.Alias {
			if strings.HasPrefix(msg, alias) {
				msg = strings.Replace(msg, fmt.Sprintf("%s ", k), "", 1)
				return msg
			}
		}

	}

	return msg
}

func createAzureWorkItem(aiResponse string) (string, error) {

	raw := aiResponse
	if !strings.HasPrefix(aiResponse, "[") {
		raw = fmt.Sprintf("[%s]", aiResponse)
	}

	var request []azure.UserStory

	err := json.Unmarshal([]byte(raw), &request)
	if err != nil {

		return "", err
	}

	workItem, err := azure.CreateAzureNewUserStory(request)
	if err != nil {
		return "", errors.New(fmt.Sprintf("ai response unexpcet %v", aiResponse))
	}
	if len(workItem.Message) > 0 {
		return "", errors.New(workItem.Message)
	}

	if len(workItem.Links.HTML.Href) == 0 {
		return "", errors.New("can not get azure url")
	}

	return fmt.Sprintf("建立 [%v] 成功 ,Url: %s", workItem.System.Title, workItem.Links.HTML.Href), nil
}

func ShowHelp(fromID string) (string, error) {

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Web Version: %s \n", "https://webgpt.sogaprod.club"))
	builder.WriteRune('\n')

	builder.WriteString(fmt.Sprintf("Chat Id: %s \n", fromID))
	builder.WriteRune('\n')

	mode := DefaultMode
	group := GetGroupMode(fromID)
	if group != nil {

		mode = group.CommandMode
	}

	builder.WriteString(fmt.Sprintf("Mode: %s \n", mode))
	builder.WriteRune('\n')

	for cmd, info := range CommandMap {

		if strings.HasPrefix(cmd, "/") {
			builder.WriteString(fmt.Sprintf("%s \n", info.Usage))
			builder.WriteRune('\n')
		}

	}

	return builder.String(), nil
}
func GetCommandBy(cmd string) *CommandInfo {

	v, ok := CommandMap[cmd]
	if !ok {
		return nil
	}
	return &v
}

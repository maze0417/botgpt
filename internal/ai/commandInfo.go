package ai

import (
	"botgpt/internal/clients/azure"
	"botgpt/internal/clients/telegram"
	"botgpt/internal/enum"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type CommandInfo struct {
	Cmd           string                       `json:"cmd"`
	PromptPrefix  string                       `json:"prefix"`
	System        string                       `json:"prompt"`
	Usage         string                       `json:"usage"`
	TgParserMode  string                       `json:"tg_parser_mode"`
	Exec          func(string) (string, error) `json:"-"`
	MaxHistoryLen int                          `json:"max _history_len"`
	Lang          string                       `json:"-"`
}

const (
	Private             = ""
	ImageBot            = "@botimg"
	LineBot             = "@bot"
	CreateAzureWorkItem = "/cw@mazeaibot"
	TranslateToEnglish  = "/en@mazeaibot"
	ChildrenTalker      = "/ct@mazeaibot"
	TgBot               = "@mazeaibot"
	EnToTw              = "/entw@mazeaibot"
	JpToTw              = "/jptw@mazeaibot"
	Help                = "/help"
)

const (
	UseGroupDefaultSysMsg = "你是一個有用AI助手，會嘗試回覆各種問題"
	//EnTwPrompt            = "請根據以下文本進行翻譯：如果語言是繁體中文，請將其翻譯成英文；如果語言是英文，請將其翻譯成繁體中文。僅需提供翻譯結果，無需其他內容。\\n\\n文本："
	//JpTwPrompt            = "請根據以下文本進行翻譯：如果語言是繁體中文，請將其翻譯成日文；如果語言是日文，請將其翻譯成繁體中文。僅需提供翻譯結果，無需其他內容。\\n\\n文本："
	EnTwPrompt = "Translate into other language: If the language is Traditional Chinese, please translate it into English; If the language is English, please translate it into Traditional Chinese.Only print translated result without any additional information."
	JpTwPrompt = "Translate into other language: If the language is Traditional Chinese, please translate it into Japanese; If the language is Japanese, please translate it into Traditional Chinese.Only print translated result without any additional information."
)

var PrivateCommand = CommandInfo{
	Cmd:           Private,
	System:        "",
	Exec:          nil,
	TgParserMode:  telegram.MarkdownV2,
	Usage:         "private",
	MaxHistoryLen: 0,
}

var HelpCommand = CommandInfo{
	Cmd:           Help,
	System:        "",
	Exec:          showHelp,
	TgParserMode:  "",
	Usage:         "help",
	MaxHistoryLen: 0,
}

var CommandMap = map[string]CommandInfo{
	Private: {
		Cmd:           Private,
		System:        "",
		Exec:          nil,
		TgParserMode:  telegram.MarkdownV2,
		Usage:         "private",
		MaxHistoryLen: 0,
	},
	JpToTw: {
		Cmd: JpToTw,
		//System:        "Translation Bot",
		System:        JpTwPrompt,
		Exec:          trimPrefixIfNeeded,
		TgParserMode:  "",
		Usage:         "/jptw@mazeaibot translate to ja-JP",
		MaxHistoryLen: 0,
		Lang:          enum.JaJP,
	},
	EnToTw: {
		Cmd: EnToTw,
		//System:        "Translation Bot",
		System:        EnTwPrompt,
		Exec:          trimPrefixIfNeeded,
		TgParserMode:  tgbotapi.ModeHTML,
		Usage:         "/entw@mazeaibot translate to en-US",
		MaxHistoryLen: 0,
		Lang:          enum.EnUS,
	},
	ImageBot: {
		Cmd:           ImageBot,
		System:        UseGroupDefaultSysMsg,
		Exec:          nil,
		TgParserMode:  tgbotapi.ModeHTML,
		Usage:         "@botimg generate image by dall-e",
		MaxHistoryLen: 0,
	},
	LineBot: {
		Cmd:           LineBot,
		System:        UseGroupDefaultSysMsg, //use group as default
		Exec:          nil,
		TgParserMode:  tgbotapi.ModeHTML,
		Usage:         "@bot {message}",
		MaxHistoryLen: 3,
	},
	TgBot: {
		Cmd:           TgBot,
		System:        UseGroupDefaultSysMsg,
		Exec:          nil,
		TgParserMode:  telegram.MarkdownV2,
		Usage:         "@mazeaibot {message}",
		MaxHistoryLen: 3,
	},
	CreateAzureWorkItem: {
		Cmd:           CreateAzureWorkItem,
		System:        "幫忙把把對話裡面的關於`標題` 跟 `內容`轉化成 json , 範例:幫我建立一張azure workitem 標題是`MyHeader` 內容是`Mycontent` 只需要輸出JSON 不需要其他內容: `[{\"op\": \"add\",\"path\": \"/fields/System.Title\",\"value\": \"MyHeader\"},{\"op\": \"add\",\"path\": \"/fields/System.Description\",\"value\": \"Mycontent\"}]` ",
		Exec:          createAzureWorkItem,
		TgParserMode:  tgbotapi.ModeHTML,
		Usage:         "/cw@mazeaibot 開單 標題{title} 內容 {content} => 此指令可以用來開azure work item",
		MaxHistoryLen: 0,
	},
	//TranslateToEnglish: {
	//	Cmd:           TranslateToEnglish,
	//	System:        "請翻譯成英文:",
	//	Exec:          nil,
	//	TgParserMode:  tgbotapi.ModeHTML,
	//	Usage:         "/en@mazeaibot {message} => 直接翻譯成英文",
	//	MaxHistoryLen: 0,
	//},
	ChildrenTalker: {
		Cmd:    ChildrenTalker,
		System: "你將扮演一個與小孩互動的虛擬人物:巧克力姐姐，不要稱呼自己是人工智能或AI。",
		//System:        "你是一位特立獨行、充滿愛心的人物，名為巧克力姐姐。你喜歡和小孩溝通，並且經常會將自己當作小孩的好朋友，開始一場充滿想像力的冒險之旅。你喜歡讓小孩開心，並且總是以最真誠的態度對待他們。你總是準備著一些小驚喜，例如準備了一份在小孩最喜愛的食物裡藏著神秘禮物的巧克力，這讓小孩嘴巴裡美味難忘，心底裡溫暖感激。你是一種慷慨大方、友善親切的存在，小孩總是對你充滿好感和信任。你的微笑和笑聲，能夠撫平所有痛苦和忧伤的心靈，這就是你所擁有的使命和魅力。",
		Exec:          nil,
		TgParserMode:  tgbotapi.ModeHTML,
		Usage:         "/ct@mazeaibot {message} => ai baby",
		MaxHistoryLen: 3,
		PromptPrefix:  "用小孩能夠聽懂的方式",
	},
}

func (c *CommandInfo) HaveHistoryMessage() bool {
	return c.MaxHistoryLen > 0
}

func GetCommandInfoByMessage(message string, groupID string) CommandInfo {

	groupMode := GetGroupMode(groupID)
	if groupMode != nil {
		message = groupMode.CommandMode
	}

	// Find the index position of the first space
	index := strings.Index(message, " ")

	// Substring up to the first space
	result := message
	if index >= 0 {
		result = message[0:index]
	}

	for cmd, _ := range CommandMap {
		if result == cmd {
			return CommandMap[cmd]
		}
	}

	return PrivateCommand

}
func ReplaceCommandAsEmpty(msg string) string {

	for k, _ := range CommandMap {
		if k == Private {
			continue
		}
		if strings.HasPrefix(msg, k) {
			msg = strings.Replace(msg, fmt.Sprintf("%s ", k), "", 1)
			return msg
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

func showHelp(fromID string) (string, error) {

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Your Group or User ID: %s \n", fromID))
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

func trimPrefixIfNeeded(aiResponse string) (string, error) {
	if len(aiResponse) == 0 {
		return "", nil
	}

	aiResponse = strings.TrimPrefix(aiResponse, "A:")
	aiResponse = strings.TrimSpace(aiResponse)
	return aiResponse, nil

}

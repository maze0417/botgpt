package ai

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/core"
	"botgpt/internal/enum"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/pkg/redis"

	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Gpt3Sender struct {
	aiProvider core.IAiProvider
}

func NewGpt3AiSender(apiProvider core.IAiProvider) core.IAiSender {
	return &Gpt3Sender{
		aiProvider: apiProvider,
	}
}

func (g Gpt3Sender) Send(messageFrom string, isGroup bool, userID string, groupID string) (error, *models.AiResponse) {
	return g.SendToGpt(messageFrom, isGroup, userID, groupID)
}

func (g Gpt3Sender) SendToGpt(messageFrom string, isGroup bool, userID string, groupID string) (error, *models.AiResponse) {

	if messageFrom == Help {
		helpText, _ := HelpCommand.Exec(groupID)
		return nil, &models.AiResponse{
			IsImage:     false,
			IsText:      true,
			Text:        helpText,
			TgParseMode: "",
			CommandMode: HelpCommand.Usage,
		}
	}

	command := GetCommandInfoByMessage(messageFrom, groupID)

	isImage := false
	isCommand := strings.HasPrefix(messageFrom, "/")
	switch command.Cmd {
	case Private:
		if isGroup {
			errMsg := fmt.Sprintf("empty Command and IsGroup , just return ")
			fmt.Println(errMsg)
			return utils.NewKnownError(enum.FALIURE, errMsg), nil
		}
	case ImageBot:
		isImage = true
	default:

	}

	if !isImage && !isCommand {
		isImage = strings.Contains(messageFrom, "draw") ||
			strings.Contains(messageFrom, "畫")
	}
	message := ReplaceCommandAsEmpty(messageFrom)

	if isCommand && len(message) <= 1 {
		return nil, &models.AiResponse{
			IsImage:     false,
			IsText:      true,
			Text:        command.Usage,
			TgParseMode: command.TgParserMode,
			CommandMode: command.Usage,
		}
	}

	if isImage {

		result, err := g.aiProvider.GenerateImage(message)

		if err != nil {
			return err, nil
		}
		log.Printf("send image with %v \n", message)

		return nil, &models.AiResponse{
			IsImage:     true,
			IsText:      false,
			Text:        result,
			TgParseMode: command.TgParserMode,
			CommandMode: command.Usage,
		}
	}
	var totalMessages []gpt3.Message
	msg := gpt3.Message{
		Role:    gpt3.User,
		Content: fmt.Sprintf("%v: %v", command.PromptPrefix, message),
	}
	err, totalMessages := getSetTotalMessages(userID, msg, command.MaxHistoryLen)
	if err != nil {
		totalMessages = append(totalMessages, msg)
	}

	if !command.HaveHistoryMessage() {
		totalMessages = []gpt3.Message{msg}
	}

	sysMsg := command.System

	totalMessages = insertSystemMessage(sysMsg, totalMessages)

	err, resp := g.aiProvider.GenerateText(totalMessages, userID)

	if err != nil {
		return err, nil
	}

	if len(resp) > 0 {
		msg := gpt3.Message{
			Role:    gpt3.Assistant,
			Content: resp,
		}
		_, _ = getSetTotalMessages(userID, msg, 1)

	}
	log.Printf("reply %v text message :: %s  \n \n ", userID, resp)
	replyToClientMsg := resp

	if err != nil {
		return err, nil
	}
	if command.Exec != nil {
		azureRes, err := command.Exec(replyToClientMsg)
		if err != nil {
			return err, nil
		}
		replyToClientMsg = azureRes
	}

	replyToClientMsg = strings.Replace(replyToClientMsg, "\n\n", "", 1)
	return nil, &models.AiResponse{
		IsImage:     false,
		IsText:      true,
		Text:        replyToClientMsg,
		TgParseMode: command.TgParserMode,
		CommandMode: command.Usage,
	}
}

func getSetTotalMessages(userID string, msg gpt3.Message, maxUserMessageLen int) (error, []gpt3.Message) {

	if maxUserMessageLen == 0 {
		return nil, nil
	}

	err, messageResult := getSetUserHistory(userID, msg, maxUserMessageLen)
	if err != nil {
		return err, nil
	}

	return nil, messageResult
}

func getSetUserHistory(userID string, msg gpt3.Message, maxUserMessageLen int) (error, []gpt3.Message) {
	// 將 msg 放進 Redis List 中，若超出最大筆數則刪除最早的一筆，使用 LRem 方法。
	b, err := json.Marshal(msg)
	if err != nil {
		return err, nil
	}

	rdb := redis.GetSingleRdb()

	ctx := context.Background()

	result, err := rdb.LLen(ctx, userID).Result()
	if err != nil {
		return err, nil
	}
	if int(result) > maxUserMessageLen {
		err = rdb.LPop(ctx, userID).Err() // 刪除最早的一筆資料
		if err != nil {
			return err, nil
		}
	}
	value := string(b)
	err = rdb.RPush(ctx, userID, value).Err()
	if err != nil {
		return err, nil
	}
	// 從 Redis List 中取出所有資料
	err, messageResult := getUserHistory(userID)
	if err != nil {
		return err, nil
	}

	return nil, messageResult
}

func getUserHistory(userID string) (error, []gpt3.Message) {

	rdb := redis.GetSingleRdb()

	ctx := context.Background()

	// 從 Redis List 中取出所有資料
	messages, err := rdb.LRange(ctx, userID, 0, -1).Result()
	if err != nil {
		return err, nil
	}

	var messageResult []gpt3.Message
	for _, m := range messages {
		var msg gpt3.Message
		err := json.Unmarshal([]byte(m), &msg)
		if err != nil {
			return fmt.Errorf("json unmarshal error: %w", err), nil
		}
		messageResult = append(messageResult, msg)
	}
	return nil, messageResult
}

//	func getAssistMessages(userID string) (error, gpt3.Message) {
//		var msg gpt3.Message
//		rdb := connection.GetSingleRdb()
//
//		ctx := context.Background()
//
//		assistKey := fmt.Sprintf("%s:assistant", userID)
//
//		lastMessage, err := rdb.Get(ctx, assistKey).Result()
//
//		if err != nil {
//			return err, msg
//		}
//
//		err = json.Unmarshal([]byte(lastMessage), &msg)
//		if err != nil {
//			return fmt.Errorf("json unmarshal error: %w", err), msg
//		}
//		return nil, msg
//	}
//
// func setAssistMessages(userID string, msg gpt3.Message) error {
//
//		b, err := json.Marshal(msg)
//		if err != nil {
//			return err
//		}
//
//		value := string(b)
//		rdb := connection.GetSingleRdb()
//
//		ctx := context.Background()
//		assistKey := fmt.Sprintf("%s:assistant", userID)
//
//		err = rdb.Set(ctx, assistKey, value, 0).Err()
//
//		return err
//	}
func insertSystemMessage(sysMsg string, totalMessages []gpt3.Message) []gpt3.Message {
	arrCopy := make([]gpt3.Message, len(totalMessages)+1)
	arrCopy[0] = gpt3.Message{
		Role:    gpt3.System,
		Content: sysMsg,
	}
	copy(arrCopy[1:], totalMessages[:])
	totalMessages = arrCopy
	return totalMessages
}

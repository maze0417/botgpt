package handler

import (
	"botgpt/internal/ai"
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/enum"
	"botgpt/internal/interfaces"
	"botgpt/internal/models"
	"botgpt/internal/utils"
	"botgpt/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

type MessageHandler struct {
	aiProvider   interfaces.IAiProvider
	textToSpeech interfaces.ITextToSpeech
}

func (m MessageHandler) SendStream(messageFrom string, userID string) (error, *openai.ChatCompletionStream) {

	msg := gpt3.Message{
		Role:    gpt3.User,
		Content: messageFrom,
	}

	var totalMessages []gpt3.Message

	err, totalMessages := getSetTotalMessages(userID, msg, 5)
	if err != nil {
		totalMessages = append(totalMessages, msg)
	}

	err, resp := m.aiProvider.GenerateTextStream(totalMessages, userID)

	if err != nil {
		return err, nil
	}

	return nil, resp

}

func NewMessageHandler(apiProvider interfaces.IAiProvider, textToSpeech interfaces.ITextToSpeech) interfaces.IMessageHandler {
	return &MessageHandler{
		aiProvider:   apiProvider,
		textToSpeech: textToSpeech,
	}
}

func (m MessageHandler) Send(messageFrom string, isGroup bool, userID string, groupID string, replyMessage string) (error, *models.AiResponse) {

	err, response, setCommand := setGroupModeIfMessageOnlyCommand(messageFrom, groupID)
	if setCommand {
		return err, response
	}

	command := ai.GetGroupModeOrCommandInfoByMessage(messageFrom, groupID)

	isImage := command.Cmd == ai.Image
	isCommand := strings.HasPrefix(messageFrom, "/")
	isAsker := command.Cmd == ai.Asker
	istrans := command.Cmd == ai.JpToTw
	ignoreChat := (command.Cmd == ai.Chat || isAsker || istrans) && isGroup && !ai.HasCommandPrefix(messageFrom)

	log.Printf("receive message from %s , group :%v , command %s ", messageFrom, isGroup, command.Cmd)

	if ignoreChat {
		errMsg := fmt.Sprintf("IsGroup and mode is chat , just return ")
		log.Println(errMsg)
		return utils.NewKnownError(enum.FALIURE, errMsg), nil
	}

	if len(replyMessage) > 0 && !isCommand {
		appendMessage(messageFrom, replyMessage)
	}

	message := ai.ReplaceCommandAsEmpty(messageFrom)

	if isCommand && len(message) <= 1 {
		return nil, &models.AiResponse{
			IsImage:     false,
			IsText:      true,
			Text:        command.Usage,
			TgParseMode: command.TgParserMode,
			CommandMode: command.Usage,
			Lang:        command.Lang,
		}
	}

	if isImage {

		result, err := m.aiProvider.GenerateImage(message, userID)

		if err != nil {
			return err, nil
		}

		return nil, &models.AiResponse{
			IsImage:     true,
			IsText:      false,
			Text:        result,
			TgParseMode: command.TgParserMode,
			CommandMode: command.Usage,
			Lang:        command.Lang,
		}
	}

	sysMsg := m.getSystemPromptFromMessage(message, command)

	llmModel := ai.GPTCheaperModel

	group := ai.GetGroupMode(groupID)
	if group != nil && group.Enable && len(group.SystemMessage) > 0 {
		sysMsg = group.SystemMessage
		if len(group.LLMModel) > 0 {
			llmModel = group.LLMModel
		}
	}
	msg := m.getPromptFromMessage(message, command)

	var totalMessages []gpt3.Message

	err, totalMessages = getSetTotalMessages(userID, msg, command.MaxHistoryLen)
	if err != nil {
		totalMessages = append(totalMessages, msg)
	}

	totalMessages = insertSystemMessage(sysMsg, totalMessages)

	if !command.HaveHistoryMessage() {
		totalMessages = []gpt3.Message{msg}
	}

	err, resp := m.aiProvider.GenerateText(totalMessages, userID, llmModel)

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

	replyToClientMsg := resp

	if err != nil {
		return err, nil
	}
	if command.PostExec != nil {
		azureRes, err := command.PostExec(replyToClientMsg)
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
		CommandMode: command.Cmd,
		Lang:        command.Lang,
	}
}

func (m MessageHandler) getPromptFromMessage(message string, cmd ai.CommandInfo) gpt3.Message {
	res := gpt3.Message{
		Role:    gpt3.User,
		Content: message,
	}

	if len(cmd.PromptPrefixInject) == 0 || len(message) == 0 {
		return res
	}
	if len(cmd.Lang) == 0 {
		return gpt3.Message{
			Role:    gpt3.User,
			Content: fmt.Sprintf("%v \n %v", cmd.PromptPrefixInject, message),
		}
	}

	targetLang := cmd.Lang

	lang := m.textToSpeech.GetLangFromText(message)

	if lang == targetLang {
		targetLang = enum.ZhTw
	}

	promptPrefixInject := fmt.Sprintf(cmd.PromptPrefixInject, targetLang)
	return gpt3.Message{
		Role:    gpt3.User,
		Content: fmt.Sprintf("%v \n %v", promptPrefixInject, message),
	}
}
func (m MessageHandler) getSystemPromptFromMessage(message string, cmd ai.CommandInfo) string {
	if len(message) == 0 {
		return cmd.SystemPrompt
	}
	if len(cmd.Lang) == 0 || len(cmd.SystemPrompt) == 0 {
		return cmd.SystemPrompt
	}

	targetLang := cmd.Lang

	lang := m.textToSpeech.GetLangFromText(message)

	if lang == targetLang {
		targetLang = enum.ZhTw
	}

	return fmt.Sprintf(cmd.SystemPrompt, targetLang)

}

func setGroupModeIfMessageOnlyCommand(messageFrom string, groupID string) (error, *models.AiResponse, bool) {
	cmd := ai.GetCommandFromAlias(messageFrom)

	if cmd == nil {
		return nil, nil, false
	}

	if cmd.Cmd == ai.Help {
		helpText, _ := ai.ShowHelp(groupID)
		return nil, &models.AiResponse{
			IsImage:     false,
			IsText:      true,
			Text:        helpText,
			TgParseMode: "",
			CommandMode: cmd.Usage,
			Lang:        cmd.Lang,
		}, true
	}

	resp := ai.SetGroupMode(groupID, cmd.Cmd)

	if resp == nil {
		return nil, nil, false
	}
	return nil, &models.AiResponse{
		IsImage:     false,
		IsText:      true,
		Text:        "mode changed to " + cmd.Cmd,
		TgParseMode: "",
		CommandMode: cmd.Usage,
		Lang:        cmd.Lang,
	}, true

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
func appendMessage(message string, append string) string {
	var builder strings.Builder
	builder.WriteString(message)
	builder.WriteRune('\n')

	builder.WriteString(append)
	builder.WriteRune('\n')

	return builder.String()
}

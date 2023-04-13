package ai

import (
	"botgpt/internal/clients/gpt3"
	"botgpt/internal/config"
	"botgpt/pkg/redis"
	"context"
	"testing"
)

type MockAiProvider struct {
}

func (v MockAiProvider) GenerateImage(message string) (string, error) {
	return "this is  mock image", nil
}

func (v MockAiProvider) GenerateText(totalMessages []gpt3.Message, userID string) (error, string) {
	return nil, "this is mock text"
}

func TestGetSetMessages(t *testing.T) {
	const maxUserMessageLen = 3
	config.InitTest("gptbot")

	rdb := redis.GetSingleRdb()
	ctx := context.Background()

	userID := "test-user"
	rdb.FlushDB(ctx)

	// Test adding messages
	msg1 := gpt3.Message{Content: "Message 1", Role: "user"}
	msg2 := gpt3.Message{Content: "Message 2", Role: "assistant"}
	msg3 := gpt3.Message{Content: "Message 3", Role: "user"}
	msg4 := gpt3.Message{Content: "Message 4", Role: "assistant"}
	msg5 := gpt3.Message{Content: "Message 5", Role: "user"}
	msg6 := gpt3.Message{Content: "Message 6", Role: "assistant"}

	err, messages := getSetTotalMessages(userID, msg1, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("Expected messages length to be 1 but got %d", len(messages))
	}
	if messages[0].Content != msg1.Content {
		t.Errorf("Expected message text to be %s but got %s", msg1.Content, messages[0].Content)
	}

	err, messages = getSetTotalMessages(userID, msg2, maxUserMessageLen)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("Expected messages length to be 2 but got %d", len(messages))
	}
	if messages[1].Content != msg2.Content {
		t.Errorf("Expected message text to be %s but got %s", msg2.Content, messages[0].Content)
	}

	err, messages = getSetTotalMessages(userID, msg3, maxUserMessageLen)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(messages) != maxUserMessageLen {
		t.Errorf("Expected messages length to be 3 but got %d", len(messages))
	}
	if messages[2].Content != msg3.Content {
		t.Errorf("Expected message text to be %s but got %s", msg3.Content, messages[0].Content)
	}

	fixLen := maxUserMessageLen + 1
	err, messages = getSetTotalMessages(userID, msg4, maxUserMessageLen)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(messages) != fixLen {
		t.Errorf("Expected messages length to be 3 but got %d", len(messages))
	}
	if messages[0].Content != msg1.Content {
		t.Errorf("Expected message text to be %s but got %s", msg1.Content, messages[0].Content)
	}

	if messages[3].Content != msg4.Content {
		t.Errorf("Expected message text to be %s but got %s", msg4.Content, messages[2].Content)
	}
	// Test removing messages
	err, messages = getSetTotalMessages(userID, msg5, maxUserMessageLen)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(messages) != fixLen {
		t.Errorf("Expected messages length to be 3 but got %d", len(messages))
	}
	if messages[0].Content != msg2.Content {
		t.Errorf("Expected message text to be %s but got %s", msg5.Content, messages[0].Content)
	}
	if messages[fixLen-1].Content != msg5.Content {
		t.Errorf("Expected message text to be %s but got %s", msg5.Content, messages[0].Content)
	}

	err, messages = getSetTotalMessages(userID, msg6, maxUserMessageLen)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(messages) != fixLen {
		t.Errorf("Expected messages length to be 3 but got %d", len(messages))
	}
	if messages[0].Content != msg3.Content {
		t.Errorf("Expected message text to be %s but got %s", msg6.Content, messages[0].Content)
	}
	if messages[fixLen-1].Content != msg6.Content {
		t.Errorf("Expected message text to be %s but got %s", msg6.Content, messages[0].Content)
	}

}

package aws

import (
	"botgpt/internal/clients/telegram"
	"botgpt/internal/config"
	"botgpt/internal/enum"
	"botgpt/internal/utils"
	"fmt"
	"github.com/aws/aws-sdk-go/service/polly"
	"testing"
)

func TestTextToSpeech(t *testing.T) {

	config.InitTest("botgpt")

	tests := []struct {
		text, lang string
	}{
		{"Hello, I am using AWS Polly with Golang.", enum.EnUS},
		{"私はチキンフィレライスを食べたいです。", enum.JaJP},
		{"你好，今天晚上天氣不錯喔", enum.CmnCN},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s+%s", tt.text, tt.lang), func(t *testing.T) {
			text := tt.text
			outputFile := utils.GetUploadDir() + "testts." + polly.OutputFormatMp3

			err := SynthesizeSpeech(text, outputFile, polly.OutputFormatMp3, tt.lang)
			if err != nil {

				t.Errorf("Failed to synthesize speech: %v\n", err)
				return
			}

			fmt.Println("Synthesized speech saved to", outputFile)

			err = telegram.SendVoice(outputFile)
			if err != nil {
				t.Error(err)
			}
		})
	}

}

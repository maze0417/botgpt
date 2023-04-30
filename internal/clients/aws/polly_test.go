package aws

import (
	"botgpt/internal/clients/telegram"
	"botgpt/internal/config"
	"botgpt/internal/enum"
	"botgpt/internal/utils"
	"fmt"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/pemistahl/lingua-go"
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

	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		Build()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s+%s", tt.text, tt.lang), func(t *testing.T) {
			text := tt.text
			var lang string
			language, exists := detector.DetectLanguageOf(text)

			if !exists {
				t.Errorf("can not detect lang from %v", text)
			}

			lang, ok := enum.LangMap[language.String()]
			if !ok {
				t.Errorf("can not map lang from %v", language.String())
			}

			outputFile := utils.GetUploadDir() + "testts." + polly.OutputFormatMp3

			err, _ := textToSpeech(text, outputFile, polly.OutputFormatMp3, lang)
			if err != nil {

				t.Errorf("Failed to synthesize speech: %v\n", err)
				return
			}

			fmt.Println("Synthesized speech saved to", outputFile)

			err = telegram.SendVoice(1066396636, outputFile)
			if err != nil {
				t.Error(err)
			}
		})
	}

}

func TestGetLangFromText(t *testing.T) {

	config.InitTest("botgpt")

	languages := []lingua.Language{
		lingua.English,
		lingua.Japanese,
		lingua.Chinese,
	}
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	texts := []string{
		"Hello, I am using AWS Polly with Golang.",
		"How old are you ?",
		"私はチキンフィレライスを食べたいです。",
		"你好，今天晚上天氣不錯喔",
	}
	// 遍歷每段文字，並識別語言
	for _, text := range texts {
		if language, exists := detector.DetectLanguageOf(text); exists {
			fmt.Printf("識別語言: %s\n", language)
		}

	}
}
